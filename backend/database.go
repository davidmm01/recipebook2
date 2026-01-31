package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const (
	localDBPath = "/tmp/recipes.db"
	dbFileName  = "recipes.db"
)

var (
	db         *sql.DB
	dbMutex    sync.RWMutex
	bucketName string
)

// Icon represents a recipe icon
type Icon struct {
	ID         string    `json:"id"`
	Filename   string    `json:"filename"`
	IconURL    string    `json:"iconUrl"`
	UploadedAt time.Time `json:"uploadedAt"`
}

// Recipe represents a recipe with markdown fields
type Recipe struct {
	ID              string        `json:"id"` // UUID
	Title           string        `json:"title"`
	Description     string        `json:"description"`     // Brief description
	RecipeType      string        `json:"type"`            // "food", "cocktail", etc.
	Cuisine         string        `json:"cuisine"`         // "italian", "japanese", "mexican", etc.
	Ingredients     string        `json:"ingredients"`     // markdown
	Method          string        `json:"method"`          // markdown
	Notes           string        `json:"notes"`           // markdown
	Sources         string        `json:"sources"`         // markdown
	IconID          *string       `json:"iconId"`          // Nullable icon ID
	Icon            *Icon         `json:"icon"`            // Icon details (loaded separately)
	Tags            []string      `json:"tags"`            // Array of tag names
	Images          []RecipeImage `json:"images"`          // Array of image URLs
	MakeCount       int           `json:"makeCount"`       // Number of times this recipe was made
	CreatedByUserID *string       `json:"createdByUserId"` // Firebase UID of creator (nullable)
	CreatedByName   *string       `json:"createdByName"`   // Display name of creator (nullable)
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

// RecipeImage represents an image associated with a recipe
type RecipeImage struct {
	ID           string    `json:"id"`
	RecipeID     string    `json:"recipeId"`
	ImageURL     string    `json:"imageUrl"`
	DisplayOrder int       `json:"displayOrder"`
	CreatedAt    time.Time `json:"createdAt"`
}

// MakeLog represents a record of when a recipe was made
type MakeLog struct {
	ID              string    `json:"id"`
	RecipeID        string    `json:"recipeId"`
	MadeAt          string    `json:"madeAt"` // Date in YYYY-MM-DD format
	Notes           string    `json:"notes"`
	CreatedByUserID *string   `json:"createdByUserId"`
	CreatedAt       time.Time `json:"createdAt"`
}

// InitDatabase initializes SQLite database and downloads from Cloud Storage if available
func InitDatabase(ctx context.Context, bucket string) error {
	bucketName = bucket

	// Try to download existing DB from Cloud Storage
	if err := downloadDBFromGCS(ctx); err != nil {
		log.Printf("No existing database found in Cloud Storage, creating new one: %v", err)
		if err := createNewDB(); err != nil {
			return fmt.Errorf("failed to create new database: %w", err)
		}
	}

	// Open SQLite connection
	var err error
	db, err = sql.Open("sqlite3", localDBPath+"?_journal_mode=WAL")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations for existing databases
	if err := runMigrations(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// createNewDB creates a new SQLite database with schema
func createNewDB() error {
	database, err := sql.Open("sqlite3", localDBPath)
	if err != nil {
		return err
	}
	defer database.Close()

	schema := `
	-- Users table (stores user profiles from Firebase Auth)
	CREATE TABLE IF NOT EXISTS users (
		firebase_uid TEXT PRIMARY KEY,
		email TEXT NOT NULL,
		display_name TEXT,
		role TEXT DEFAULT 'viewer' CHECK(role IN ('viewer', 'editor', 'admin')),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_login_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

	-- Icons table (shared collection of recipe icons)
	CREATE TABLE IF NOT EXISTS icons (
		id TEXT PRIMARY KEY,
		filename TEXT NOT NULL,
		icon_url TEXT NOT NULL,
		uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS recipes (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		recipe_type TEXT,
		cuisine TEXT,
		ingredients TEXT,
		method TEXT,
		notes TEXT,
		sources TEXT,
		icon_id TEXT,
		created_by_user_id TEXT,
		created_by_name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (icon_id) REFERENCES icons(id)
	);

	CREATE INDEX IF NOT EXISTS idx_recipes_type ON recipes(recipe_type);
	CREATE INDEX IF NOT EXISTS idx_recipes_cuisine ON recipes(cuisine);
	CREATE INDEX IF NOT EXISTS idx_recipes_icon ON recipes(icon_id);

	-- Tags table
	CREATE TABLE IF NOT EXISTS tags (
		id TEXT PRIMARY KEY,
		name TEXT UNIQUE NOT NULL
	);

	-- Many-to-many relationship between recipes and tags
	CREATE TABLE IF NOT EXISTS recipe_tags (
		recipe_id TEXT NOT NULL,
		tag_id TEXT NOT NULL,
		PRIMARY KEY (recipe_id, tag_id),
		FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_recipe_tags_recipe ON recipe_tags(recipe_id);
	CREATE INDEX IF NOT EXISTS idx_recipe_tags_tag ON recipe_tags(tag_id);

	-- Recipe images table
	CREATE TABLE IF NOT EXISTS recipe_images (
		id TEXT PRIMARY KEY,
		recipe_id TEXT NOT NULL,
		image_url TEXT NOT NULL,
		display_order INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_recipe_images_recipe ON recipe_images(recipe_id);

	-- Make logs table (tracks when recipes were made)
	CREATE TABLE IF NOT EXISTS make_logs (
		id TEXT PRIMARY KEY,
		recipe_id TEXT NOT NULL,
		made_at DATE NOT NULL,
		notes TEXT,
		created_by_user_id TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
		FOREIGN KEY (created_by_user_id) REFERENCES users(firebase_uid)
	);

	CREATE INDEX IF NOT EXISTS idx_make_logs_recipe ON make_logs(recipe_id);
	CREATE INDEX IF NOT EXISTS idx_make_logs_made_at ON make_logs(made_at);

	-- Full-text search table (now includes description, cuisine, and sources)
	-- Note: We include recipe_id as an unindexed column to enable joining back to recipes table
	CREATE VIRTUAL TABLE IF NOT EXISTS recipes_fts USING fts5(
		recipe_id UNINDEXED,
		title, description, cuisine, ingredients, method, notes, sources
	);

	-- Triggers to keep FTS in sync
	CREATE TRIGGER IF NOT EXISTS recipes_fts_insert AFTER INSERT ON recipes BEGIN
		INSERT INTO recipes_fts(recipe_id, title, description, cuisine, ingredients, method, notes, sources)
		VALUES (new.id, new.title, new.description, new.cuisine, new.ingredients, new.method, new.notes, new.sources);
	END;

	CREATE TRIGGER IF NOT EXISTS recipes_fts_update AFTER UPDATE ON recipes BEGIN
		UPDATE recipes_fts
		SET title=new.title, description=new.description, cuisine=new.cuisine,
			ingredients=new.ingredients, method=new.method, notes=new.notes, sources=new.sources
		WHERE recipe_id=new.id;
	END;

	CREATE TRIGGER IF NOT EXISTS recipes_fts_delete AFTER DELETE ON recipes BEGIN
		DELETE FROM recipes_fts WHERE recipe_id=old.id;
	END;
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("Created new database with schema")
	return nil
}

// runMigrations applies database migrations for schema changes
func runMigrations(ctx context.Context) error {
	// Migration 1: Add creator fields to recipes table
	var columnExists int
	err := db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('recipes') WHERE name='created_by_user_id'").Scan(&columnExists)
	if err != nil {
		return fmt.Errorf("failed to check for column existence: %w", err)
	}

	if columnExists == 0 {
		log.Println("Running migration: Adding creator fields to recipes table")
		migrations := []string{
			"ALTER TABLE recipes ADD COLUMN created_by_user_id TEXT",
			"ALTER TABLE recipes ADD COLUMN created_by_name TEXT",
		}

		for _, migration := range migrations {
			if _, err := db.ExecContext(ctx, migration); err != nil {
				return fmt.Errorf("failed to run migration '%s': %w", migration, err)
			}
		}
		log.Println("Migration completed: Creator fields added")
	}

	// Migration 2: Add users table
	var tableExists int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check for users table existence: %w", err)
	}

	if tableExists == 0 {
		log.Println("Running migration: Creating users table")
		userTableSQL := `
			CREATE TABLE users (
				firebase_uid TEXT PRIMARY KEY,
				email TEXT NOT NULL,
				display_name TEXT,
				role TEXT DEFAULT 'viewer' CHECK(role IN ('viewer', 'editor', 'admin')),
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				last_login_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);
			CREATE INDEX idx_users_email ON users(email);
		`
		if _, err := db.ExecContext(ctx, userTableSQL); err != nil {
			return fmt.Errorf("failed to create users table: %w", err)
		}
		log.Println("Migration completed: Users table created")
	}

	// Migration 3: Add make_logs table
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='make_logs'").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check for make_logs table existence: %w", err)
	}

	if tableExists == 0 {
		log.Println("Running migration: Creating make_logs table")
		makeLogsTableSQL := `
			CREATE TABLE make_logs (
				id TEXT PRIMARY KEY,
				recipe_id TEXT NOT NULL,
				made_at DATE NOT NULL,
				notes TEXT,
				created_by_user_id TEXT,
				created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
				FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
				FOREIGN KEY (created_by_user_id) REFERENCES users(firebase_uid)
			);
			CREATE INDEX idx_make_logs_recipe ON make_logs(recipe_id);
			CREATE INDEX idx_make_logs_made_at ON make_logs(made_at);
		`
		if _, err := db.ExecContext(ctx, makeLogsTableSQL); err != nil {
			return fmt.Errorf("failed to create make_logs table: %w", err)
		}
		log.Println("Migration completed: make_logs table created")
	}

	// Migration 4: Convert integer IDs to UUID strings
	if err := migrateIntegerIDsToUUIDs(ctx); err != nil {
		return fmt.Errorf("failed to migrate integer IDs to UUIDs: %w", err)
	}

	return nil
}

// migrateIntegerIDsToUUIDs migrates all integer primary keys to UUID strings
func migrateIntegerIDsToUUIDs(ctx context.Context) error {
	// Check if icons table uses INTEGER ID (if so, migration needed)
	var iconIDType string
	err := db.QueryRow(`SELECT type FROM pragma_table_info('icons') WHERE name='id'`).Scan(&iconIDType)
	if err != nil {
		// Table might not exist yet
		return nil
	}

	// If id is already TEXT, migration already completed
	if iconIDType == "TEXT" {
		return nil
	}

	log.Println("Running migration: Converting integer IDs to UUIDs")

	// Start a transaction for the migration
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Migrate icons table
	log.Println("  - Migrating icons table")
	iconIDMap := make(map[int64]string)

	// Get all existing icons
	rows, err := tx.QueryContext(ctx, "SELECT id, filename, icon_url, uploaded_at FROM icons")
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to query icons: %w", err)
	}
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var oldID int64
			var filename, iconURL string
			var uploadedAt time.Time
			if err := rows.Scan(&oldID, &filename, &iconURL, &uploadedAt); err != nil {
				return fmt.Errorf("failed to scan icon: %w", err)
			}
			iconIDMap[oldID] = uuid.New().String()
		}
		rows.Close()
	}

	// Create new icons table with UUID
	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE icons_new (
			id TEXT PRIMARY KEY,
			filename TEXT NOT NULL,
			icon_url TEXT NOT NULL,
			uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`); err != nil {
		return fmt.Errorf("failed to create icons_new table: %w", err)
	}

	// Copy data with new UUIDs
	for oldID, newID := range iconIDMap {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO icons_new (id, filename, icon_url, uploaded_at)
			SELECT ?, filename, icon_url, uploaded_at FROM icons WHERE id = ?
		`, newID, oldID)
		if err != nil {
			return fmt.Errorf("failed to copy icon data: %w", err)
		}
	}

	// 2. Migrate tags table
	log.Println("  - Migrating tags table")
	tagIDMap := make(map[int64]string)

	rows, err = tx.QueryContext(ctx, "SELECT id, name FROM tags")
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to query tags: %w", err)
	}
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var oldID int64
			var name string
			if err := rows.Scan(&oldID, &name); err != nil {
				return fmt.Errorf("failed to scan tag: %w", err)
			}
			tagIDMap[oldID] = uuid.New().String()
		}
		rows.Close()
	}

	// Create new tags table
	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE tags_new (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL
		)
	`); err != nil {
		return fmt.Errorf("failed to create tags_new table: %w", err)
	}

	// Copy data with new UUIDs
	for oldID, newID := range tagIDMap {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO tags_new (id, name)
			SELECT ?, name FROM tags WHERE id = ?
		`, newID, oldID)
		if err != nil {
			return fmt.Errorf("failed to copy tag data: %w", err)
		}
	}

	// 3. Update recipes table to use TEXT for icon_id
	log.Println("  - Migrating recipes.icon_id")
	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE recipes_new (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			recipe_type TEXT,
			cuisine TEXT,
			ingredients TEXT,
			method TEXT,
			notes TEXT,
			sources TEXT,
			icon_id TEXT,
			created_by_user_id TEXT,
			created_by_name TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (icon_id) REFERENCES icons(id)
		)
	`); err != nil {
		return fmt.Errorf("failed to create recipes_new table: %w", err)
	}

	// Copy recipes and update icon_id references
	recipeRows, err := tx.QueryContext(ctx, "SELECT id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, icon_id, created_by_user_id, created_by_name, created_at, updated_at FROM recipes")
	if err != nil {
		return fmt.Errorf("failed to query recipes: %w", err)
	}
	defer recipeRows.Close()

	for recipeRows.Next() {
		var id, title, description, recipeType, cuisine, ingredients, method, notes, sources string
		var iconID sql.NullInt64
		var createdByUserID, createdByName sql.NullString
		var createdAt, updatedAt time.Time

		if err := recipeRows.Scan(&id, &title, &description, &recipeType, &cuisine, &ingredients, &method, &notes, &sources, &iconID, &createdByUserID, &createdByName, &createdAt, &updatedAt); err != nil {
			return fmt.Errorf("failed to scan recipe: %w", err)
		}

		var newIconID sql.NullString
		if iconID.Valid {
			if mappedID, ok := iconIDMap[iconID.Int64]; ok {
				newIconID = sql.NullString{String: mappedID, Valid: true}
			}
		}

		_, err := tx.ExecContext(ctx, `
			INSERT INTO recipes_new (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, icon_id, created_by_user_id, created_by_name, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, id, title, description, recipeType, cuisine, ingredients, method, notes, sources, newIconID, createdByUserID, createdByName, createdAt, updatedAt)
		if err != nil {
			return fmt.Errorf("failed to copy recipe: %w", err)
		}
	}

	// 4. Migrate recipe_tags table
	log.Println("  - Migrating recipe_tags table")
	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE recipe_tags_new (
			recipe_id TEXT NOT NULL,
			tag_id TEXT NOT NULL,
			PRIMARY KEY (recipe_id, tag_id),
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		)
	`); err != nil {
		return fmt.Errorf("failed to create recipe_tags_new table: %w", err)
	}

	// Copy recipe_tags with updated tag IDs
	rtRows, err := tx.QueryContext(ctx, "SELECT recipe_id, tag_id FROM recipe_tags")
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to query recipe_tags: %w", err)
	}
	if rtRows != nil {
		defer rtRows.Close()
		for rtRows.Next() {
			var recipeID string
			var oldTagID int64
			if err := rtRows.Scan(&recipeID, &oldTagID); err != nil {
				return fmt.Errorf("failed to scan recipe_tag: %w", err)
			}

			if newTagID, ok := tagIDMap[oldTagID]; ok {
				_, err := tx.ExecContext(ctx, `
					INSERT INTO recipe_tags_new (recipe_id, tag_id) VALUES (?, ?)
				`, recipeID, newTagID)
				if err != nil {
					return fmt.Errorf("failed to copy recipe_tag: %w", err)
				}
			}
		}
	}

	// 5. Migrate recipe_images table
	log.Println("  - Migrating recipe_images table")
	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE recipe_images_new (
			id TEXT PRIMARY KEY,
			recipe_id TEXT NOT NULL,
			image_url TEXT NOT NULL,
			display_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
		)
	`); err != nil {
		return fmt.Errorf("failed to create recipe_images_new table: %w", err)
	}

	riRows, err := tx.QueryContext(ctx, "SELECT id, recipe_id, image_url, display_order, created_at FROM recipe_images")
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to query recipe_images: %w", err)
	}
	if riRows != nil {
		defer riRows.Close()
		for riRows.Next() {
			var oldID int64
			var recipeID, imageURL string
			var displayOrder int
			var createdAt time.Time
			if err := riRows.Scan(&oldID, &recipeID, &imageURL, &displayOrder, &createdAt); err != nil {
				return fmt.Errorf("failed to scan recipe_image: %w", err)
			}

			newID := uuid.New().String()
			_, err := tx.ExecContext(ctx, `
				INSERT INTO recipe_images_new (id, recipe_id, image_url, display_order, created_at)
				VALUES (?, ?, ?, ?, ?)
			`, newID, recipeID, imageURL, displayOrder, createdAt)
			if err != nil {
				return fmt.Errorf("failed to copy recipe_image: %w", err)
			}
		}
	}

	// 6. Migrate make_logs table
	log.Println("  - Migrating make_logs table")
	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE make_logs_new (
			id TEXT PRIMARY KEY,
			recipe_id TEXT NOT NULL,
			made_at DATE NOT NULL,
			notes TEXT,
			created_by_user_id TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
			FOREIGN KEY (created_by_user_id) REFERENCES users(firebase_uid)
		)
	`); err != nil {
		return fmt.Errorf("failed to create make_logs_new table: %w", err)
	}

	mlRows, err := tx.QueryContext(ctx, "SELECT id, recipe_id, made_at, notes, created_by_user_id, created_at FROM make_logs")
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to query make_logs: %w", err)
	}
	if mlRows != nil {
		defer mlRows.Close()
		for mlRows.Next() {
			var oldID int64
			var recipeID, madeAt string
			var notes, createdByUserID sql.NullString
			var createdAt time.Time
			if err := mlRows.Scan(&oldID, &recipeID, &madeAt, &notes, &createdByUserID, &createdAt); err != nil {
				return fmt.Errorf("failed to scan make_log: %w", err)
			}

			newID := uuid.New().String()
			_, err := tx.ExecContext(ctx, `
				INSERT INTO make_logs_new (id, recipe_id, made_at, notes, created_by_user_id, created_at)
				VALUES (?, ?, ?, ?, ?, ?)
			`, newID, recipeID, madeAt, notes, createdByUserID, createdAt)
			if err != nil {
				return fmt.Errorf("failed to copy make_log: %w", err)
			}
		}
	}

	// Drop old tables and rename new ones
	log.Println("  - Swapping tables")

	// Drop FTS triggers first
	if _, err := tx.ExecContext(ctx, "DROP TRIGGER IF EXISTS recipes_fts_insert"); err != nil {
		return fmt.Errorf("failed to drop trigger: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DROP TRIGGER IF EXISTS recipes_fts_update"); err != nil {
		return fmt.Errorf("failed to drop trigger: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DROP TRIGGER IF EXISTS recipes_fts_delete"); err != nil {
		return fmt.Errorf("failed to drop trigger: %w", err)
	}

	// Drop old tables
	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS recipe_tags"); err != nil {
		return fmt.Errorf("failed to drop recipe_tags: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS recipe_images"); err != nil {
		return fmt.Errorf("failed to drop recipe_images: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS make_logs"); err != nil {
		return fmt.Errorf("failed to drop make_logs: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS recipes"); err != nil {
		return fmt.Errorf("failed to drop recipes: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS tags"); err != nil {
		return fmt.Errorf("failed to drop tags: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS icons"); err != nil {
		return fmt.Errorf("failed to drop icons: %w", err)
	}

	// Rename new tables
	if _, err := tx.ExecContext(ctx, "ALTER TABLE icons_new RENAME TO icons"); err != nil {
		return fmt.Errorf("failed to rename icons_new: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "ALTER TABLE tags_new RENAME TO tags"); err != nil {
		return fmt.Errorf("failed to rename tags_new: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "ALTER TABLE recipes_new RENAME TO recipes"); err != nil {
		return fmt.Errorf("failed to rename recipes_new: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "ALTER TABLE recipe_tags_new RENAME TO recipe_tags"); err != nil {
		return fmt.Errorf("failed to rename recipe_tags_new: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "ALTER TABLE recipe_images_new RENAME TO recipe_images"); err != nil {
		return fmt.Errorf("failed to rename recipe_images_new: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "ALTER TABLE make_logs_new RENAME TO make_logs"); err != nil {
		return fmt.Errorf("failed to rename make_logs_new: %w", err)
	}

	// Recreate indexes
	log.Println("  - Recreating indexes")
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_recipes_type ON recipes(recipe_type)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_cuisine ON recipes(cuisine)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_icon ON recipes(icon_id)",
		"CREATE INDEX IF NOT EXISTS idx_recipe_tags_recipe ON recipe_tags(recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_recipe_tags_tag ON recipe_tags(tag_id)",
		"CREATE INDEX IF NOT EXISTS idx_recipe_images_recipe ON recipe_images(recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_make_logs_recipe ON make_logs(recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_make_logs_made_at ON make_logs(made_at)",
	}
	for _, idx := range indexes {
		if _, err := tx.ExecContext(ctx, idx); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	// Recreate FTS triggers
	log.Println("  - Recreating FTS triggers")
	if _, err := tx.ExecContext(ctx, `
		CREATE TRIGGER recipes_fts_insert AFTER INSERT ON recipes BEGIN
			INSERT INTO recipes_fts(recipe_id, title, description, cuisine, ingredients, method, notes, sources)
			VALUES (new.id, new.title, new.description, new.cuisine, new.ingredients, new.method, new.notes, new.sources);
		END
	`); err != nil {
		return fmt.Errorf("failed to create insert trigger: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		CREATE TRIGGER recipes_fts_update AFTER UPDATE ON recipes BEGIN
			UPDATE recipes_fts
			SET title=new.title, description=new.description, cuisine=new.cuisine,
				ingredients=new.ingredients, method=new.method, notes=new.notes, sources=new.sources
			WHERE recipe_id=new.id;
		END
	`); err != nil {
		return fmt.Errorf("failed to create update trigger: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		CREATE TRIGGER recipes_fts_delete AFTER DELETE ON recipes BEGIN
			DELETE FROM recipes_fts WHERE recipe_id=old.id;
		END
	`); err != nil {
		return fmt.Errorf("failed to create delete trigger: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration transaction: %w", err)
	}

	log.Println("Migration completed: Converted all integer IDs to UUIDs")
	return nil
}

// downloadDBFromGCS downloads the SQLite database from Cloud Storage
func downloadDBFromGCS(ctx context.Context) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	rc, err := client.Bucket(bucketName).Object(dbFileName).NewReader(ctx)
	if err != nil {
		return err
	}
	defer rc.Close()

	f, err := os.Create(localDBPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = io.Copy(f, rc); err != nil {
		return err
	}

	log.Println("Downloaded database from Cloud Storage")
	return nil
}

// uploadDBToGCS uploads the SQLite database to Cloud Storage.
// Callers must hold dbMutex.
//
// This is called synchronously during write operations, which adds latency to
// every mutation request. We can't do this asynchronously in a goroutine because
// Cloud Run throttles CPU when there are no active requests (default behaviour),
// so background goroutines get frozen before the upload completes. The local
// SQLite file lives in /tmp and is lost when the instance is recycled, meaning
// any writes that weren't uploaded to GCS are silently dropped. Synchronous
// uploads avoid this at the cost of slower write responses.
//
// The alternative is setting Cloud Run to "CPU always allocated"
// (--no-cpu-throttling), which would let async goroutines finish, but that
// increases the bill since you pay for CPU even between requests. Not worth it
// at this stage.
func uploadDBToGCS(ctx context.Context) error {
	// Force WAL checkpoint so all writes are flushed to the main .db file.
	// Without this, recent writes may only exist in the -wal file and won't
	// be included when we read/upload the main database file.
	if _, err := db.ExecContext(ctx, "PRAGMA wal_checkpoint(TRUNCATE)"); err != nil {
		return fmt.Errorf("failed to checkpoint WAL: %w", err)
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	f, err := os.Open(localDBPath)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := client.Bucket(bucketName).Object(dbFileName).NewWriter(ctx)
	wc.ContentType = "application/x-sqlite3"

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	log.Println("Uploaded database to Cloud Storage")
	return nil
}

// GetRecipes returns all recipes
func GetRecipes(ctx context.Context) ([]Recipe, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	query := `
		SELECT id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, icon_id, created_by_user_id, created_by_name, created_at, updated_at
		FROM recipes
		ORDER BY updated_at DESC
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.Sources, &r.IconID, &r.CreatedByUserID, &r.CreatedByName, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Load icon if iconID is set
		if r.IconID != nil {
			icon, err := getIconByID(ctx, *r.IconID)
			if err == nil {
				r.Icon = icon
			}
		}

		// Load tags for this recipe
		tags, err := getRecipeTags(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.Tags = tags

		// Load images for this recipe
		images, err := getRecipeImages(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.Images = images

		// Load make count for this recipe
		makeCount, err := GetMakeCountByRecipe(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.MakeCount = makeCount

		recipes = append(recipes, r)
	}

	return recipes, rows.Err()
}

// GetRecipeByID returns a single recipe by ID
func GetRecipeByID(ctx context.Context, recipeID string) (*Recipe, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	query := `
		SELECT id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, icon_id, created_by_user_id, created_by_name, created_at, updated_at
		FROM recipes
		WHERE id = ?
	`

	var r Recipe
	err := db.QueryRowContext(ctx, query, recipeID).Scan(
		&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.Sources, &r.IconID, &r.CreatedByUserID, &r.CreatedByName, &r.CreatedAt, &r.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Load icon if iconID is set
	if r.IconID != nil {
		icon, err := getIconByID(ctx, *r.IconID)
		if err == nil {
			r.Icon = icon
		}
	}

	// Load tags for this recipe
	tags, err := getRecipeTags(ctx, r.ID)
	if err != nil {
		return nil, err
	}
	r.Tags = tags

	// Load images for this recipe
	images, err := getRecipeImages(ctx, r.ID)
	if err != nil {
		return nil, err
	}
	r.Images = images

	// Load make count for this recipe
	makeCount, err := GetMakeCountByRecipe(ctx, r.ID)
	if err != nil {
		return nil, err
	}
	r.MakeCount = makeCount

	return &r, nil
}

// prepareFTS5Query prepares a search query for FTS5 prefix matching
// Converts "bou" to "bou*" and "lime juice" to "lime* juice*" for better partial matching
func prepareFTS5Query(query string) string {
	// Split by whitespace
	words := strings.Fields(query)
	if len(words) == 0 {
		return query
	}

	// Add prefix matching (*) to each word that doesn't already have FTS5 operators
	var prepared []string
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		// Don't add * if the word already has FTS5 operators like *, OR, AND, NOT
		if !strings.Contains(word, "*") &&
			!strings.EqualFold(word, "OR") &&
			!strings.EqualFold(word, "AND") &&
			!strings.EqualFold(word, "NOT") {
			word = word + "*"
		}
		prepared = append(prepared, word)
	}

	return strings.Join(prepared, " ")
}

// SearchRecipes performs full-text search across recipes
func SearchRecipes(ctx context.Context, query string) ([]Recipe, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	sqlQuery := `
		SELECT r.id, r.title, r.description, r.recipe_type, r.cuisine, r.ingredients, r.method, r.notes, r.sources, r.icon_id, r.created_by_user_id, r.created_by_name, r.created_at, r.updated_at
		FROM recipes r
		JOIN recipes_fts ON r.id = recipes_fts.recipe_id
		WHERE recipes_fts MATCH ?
		ORDER BY rank
	`

	// Prepare query for prefix matching
	preparedQuery := prepareFTS5Query(query)
	rows, err := db.QueryContext(ctx, sqlQuery, preparedQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.Sources, &r.IconID, &r.CreatedByUserID, &r.CreatedByName, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Load icon if iconID is set
		if r.IconID != nil {
			icon, err := getIconByID(ctx, *r.IconID)
			if err == nil {
				r.Icon = icon
			}
		}

		// Load tags for this recipe
		tags, err := getRecipeTags(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.Tags = tags

		// Load images for this recipe
		images, err := getRecipeImages(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.Images = images

		// Load make count for this recipe
		makeCount, err := GetMakeCountByRecipe(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.MakeCount = makeCount

		recipes = append(recipes, r)
	}

	return recipes, rows.Err()
}

// FilterRecipes performs filtering and sorting based on search text, tags, cuisine, recipe type, and sort order
// If searchQuery is provided, it uses FTS5 for text search
// If tags are provided, filters recipes that have ALL specified tags
// If cuisine is provided, filters by exact cuisine match
// If recipeType is provided, filters by recipe type (food, drink, etc.)
// If sortBy is provided, sorts results accordingly
func FilterRecipes(ctx context.Context, searchQuery string, tags []string, cuisine string, recipeType string, sortBy string) ([]Recipe, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	var queryBuilder strings.Builder
	var args []interface{}

	// Base query - start with recipes table
	if searchQuery != "" {
		// Use FTS5 for text search with prefix matching
		queryBuilder.WriteString(`
			SELECT DISTINCT r.id, r.title, r.description, r.recipe_type, r.cuisine, r.ingredients, r.method, r.notes, r.sources, r.icon_id, r.created_by_user_id, r.created_by_name, r.created_at, r.updated_at
			FROM recipes r
			JOIN recipes_fts ON r.id = recipes_fts.recipe_id
			WHERE recipes_fts MATCH ?
		`)
		args = append(args, prepareFTS5Query(searchQuery))
	} else {
		// No text search, just filter
		queryBuilder.WriteString(`
			SELECT DISTINCT r.id, r.title, r.description, r.recipe_type, r.cuisine, r.ingredients, r.method, r.notes, r.sources, r.icon_id, r.created_by_user_id, r.created_by_name, r.created_at, r.updated_at
			FROM recipes r
			WHERE 1=1
		`)
	}

	// Add recipe type filter
	if recipeType != "" {
		queryBuilder.WriteString(` AND r.recipe_type = ?`)
		args = append(args, recipeType)
	}

	// Add cuisine filter
	if cuisine != "" {
		queryBuilder.WriteString(` AND r.cuisine = ?`)
		args = append(args, cuisine)
	}

	// Add tag filters - recipe must have ALL specified tags
	if len(tags) > 0 {
		queryBuilder.WriteString(`
			AND r.id IN (
				SELECT rt.recipe_id
				FROM recipe_tags rt
				JOIN tags t ON rt.tag_id = t.id
				WHERE t.name IN (`)

		for i, tag := range tags {
			if i > 0 {
				queryBuilder.WriteString(`, `)
			}
			queryBuilder.WriteString(`?`)
			args = append(args, strings.TrimSpace(strings.ToLower(tag)))
		}

		queryBuilder.WriteString(`)
				GROUP BY rt.recipe_id
				HAVING COUNT(DISTINCT t.id) = ?
			)`)
		args = append(args, len(tags))
	}

	// Order by
	if searchQuery != "" && sortBy == "" {
		// Default to rank sorting when searching
		queryBuilder.WriteString(` ORDER BY rank`)
	} else {
		// Apply sorting based on sortBy parameter
		switch sortBy {
		case "created_desc":
			queryBuilder.WriteString(` ORDER BY r.created_at DESC`)
		case "created_asc":
			queryBuilder.WriteString(` ORDER BY r.created_at ASC`)
		case "name_asc":
			queryBuilder.WriteString(` ORDER BY r.title COLLATE NOCASE ASC`)
		case "name_desc":
			queryBuilder.WriteString(` ORDER BY r.title COLLATE NOCASE DESC`)
		case "made_desc":
			queryBuilder.WriteString(` ORDER BY (SELECT COUNT(*) FROM make_logs WHERE recipe_id = r.id) DESC`)
		case "made_asc":
			queryBuilder.WriteString(` ORDER BY (SELECT COUNT(*) FROM make_logs WHERE recipe_id = r.id) ASC`)
		case "updated_desc":
			fallthrough
		default:
			queryBuilder.WriteString(` ORDER BY r.updated_at DESC`)
		}
	}

	rows, err := db.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.Sources, &r.IconID, &r.CreatedByUserID, &r.CreatedByName, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Load icon if iconID is set
		if r.IconID != nil {
			icon, err := getIconByID(ctx, *r.IconID)
			if err == nil {
				r.Icon = icon
			}
		}

		// Load tags for this recipe
		recipeTags, err := getRecipeTags(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.Tags = recipeTags

		// Load images for this recipe
		images, err := getRecipeImages(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.Images = images

		// Load make count for this recipe
		makeCount, err := GetMakeCountByRecipe(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		r.MakeCount = makeCount

		recipes = append(recipes, r)
	}

	return recipes, rows.Err()
}

// CreateRecipe inserts a new recipe and syncs to Cloud Storage
func CreateRecipe(ctx context.Context, recipe *Recipe) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Generate UUID if not provided
	if recipe.ID == "" {
		recipe.ID = uuid.New().String()
	}

	recipe.CreatedAt = time.Now()
	recipe.UpdatedAt = time.Now()

	query := `
		INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, icon_id, created_by_user_id, created_by_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.ExecContext(ctx, query,
		recipe.ID, recipe.Title, recipe.Description, recipe.RecipeType, recipe.Cuisine,
		recipe.Ingredients, recipe.Method, recipe.Notes, recipe.Sources, recipe.IconID,
		recipe.CreatedByUserID, recipe.CreatedByName, recipe.CreatedAt, recipe.UpdatedAt,
	)
	if err != nil {
		return err
	}

	// Handle tags
	if len(recipe.Tags) > 0 {
		if err := setRecipeTags(ctx, recipe.ID, recipe.Tags); err != nil {
			return err
		}
	}

	// Upload to Cloud Storage
	if err := uploadDBToGCS(context.Background()); err != nil {
		log.Printf("Failed to upload database to Cloud Storage: %v", err)
	}

	return nil
}

// UpdateRecipe updates an existing recipe and syncs to Cloud Storage
func UpdateRecipe(ctx context.Context, recipe *Recipe) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	recipe.UpdatedAt = time.Now()

	query := `
		UPDATE recipes
		SET title = ?, description = ?, recipe_type = ?, cuisine = ?,
		    ingredients = ?, method = ?, notes = ?, sources = ?, icon_id = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := db.ExecContext(ctx, query,
		recipe.Title, recipe.Description, recipe.RecipeType, recipe.Cuisine,
		recipe.Ingredients, recipe.Method, recipe.Notes, recipe.Sources, recipe.IconID, recipe.UpdatedAt,
		recipe.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("recipe not found")
	}

	// Update tags (remove old ones and add new ones)
	if err := setRecipeTags(ctx, recipe.ID, recipe.Tags); err != nil {
		return err
	}

	// Upload to Cloud Storage
	if err := uploadDBToGCS(context.Background()); err != nil {
		log.Printf("Failed to upload database to Cloud Storage: %v", err)
	}

	return nil
}

// DeleteRecipe deletes a recipe and syncs to Cloud Storage
func DeleteRecipe(ctx context.Context, recipeID string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	query := `DELETE FROM recipes WHERE id = ?`
	result, err := db.ExecContext(ctx, query, recipeID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("recipe not found")
	}

	// Tags will be automatically deleted via ON DELETE CASCADE

	// Upload to Cloud Storage
	if err := uploadDBToGCS(context.Background()); err != nil {
		log.Printf("Failed to upload database to Cloud Storage: %v", err)
	}

	return nil
}

// Helper functions for icon management

// getIconByID returns an icon by ID
func getIconByID(ctx context.Context, iconID string) (*Icon, error) {
	query := `SELECT id, filename, icon_url, uploaded_at FROM icons WHERE id = ?`
	var icon Icon
	err := db.QueryRowContext(ctx, query, iconID).Scan(&icon.ID, &icon.Filename, &icon.IconURL, &icon.UploadedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &icon, nil
}

// GetAllIcons returns all icons
func GetAllIcons(ctx context.Context) ([]Icon, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	query := `SELECT id, filename, icon_url, uploaded_at FROM icons ORDER BY uploaded_at DESC`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var icons []Icon
	for rows.Next() {
		var icon Icon
		if err := rows.Scan(&icon.ID, &icon.Filename, &icon.IconURL, &icon.UploadedAt); err != nil {
			return nil, err
		}
		icons = append(icons, icon)
	}

	return icons, rows.Err()
}

// createIcon creates a new icon record
func createIcon(ctx context.Context, filename, iconURL string) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO icons (id, filename, icon_url) VALUES (?, ?, ?)`
	_, err := db.ExecContext(ctx, query, id, filename, iconURL)
	if err != nil {
		return "", err
	}
	return id, nil
}

// GetAllTags returns all unique tags in the database, optionally filtered by recipe type
func GetAllTags(ctx context.Context, recipeType string) ([]string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	var query string
	var args []interface{}

	if recipeType != "" {
		// Filter tags by recipe type
		query = `
			SELECT DISTINCT t.name
			FROM tags t
			JOIN recipe_tags rt ON t.id = rt.tag_id
			JOIN recipes r ON rt.recipe_id = r.id
			WHERE r.recipe_type = ?
			ORDER BY t.name
		`
		args = append(args, recipeType)
	} else {
		// Get all tags
		query = `SELECT name FROM tags ORDER BY name`
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// GetAllCuisines returns all unique cuisines in the database, optionally filtered by recipe type
func GetAllCuisines(ctx context.Context, recipeType string) ([]string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	var query string
	var args []interface{}

	if recipeType != "" {
		// Filter cuisines by recipe type
		query = `
			SELECT DISTINCT cuisine
			FROM recipes
			WHERE cuisine IS NOT NULL AND cuisine != '' AND recipe_type = ?
			ORDER BY cuisine
		`
		args = append(args, recipeType)
	} else {
		// Get all cuisines
		query = `SELECT DISTINCT cuisine FROM recipes WHERE cuisine IS NOT NULL AND cuisine != '' ORDER BY cuisine`
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cuisines []string
	for rows.Next() {
		var cuisine string
		if err := rows.Scan(&cuisine); err != nil {
			return nil, err
		}
		cuisines = append(cuisines, cuisine)
	}

	return cuisines, rows.Err()
}

// Helper functions for tag management

// getRecipeTags returns all tags for a given recipe
func getRecipeTags(ctx context.Context, recipeID string) ([]string, error) {
	query := `
		SELECT t.name
		FROM tags t
		JOIN recipe_tags rt ON t.id = rt.tag_id
		WHERE rt.recipe_id = ?
		ORDER BY t.name
	`

	rows, err := db.QueryContext(ctx, query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// setRecipeTags updates the tags for a recipe (removes old ones, adds new ones)
func setRecipeTags(ctx context.Context, recipeID string, tagNames []string) error {
	// Remove existing tag associations
	deleteQuery := `DELETE FROM recipe_tags WHERE recipe_id = ?`
	if _, err := db.ExecContext(ctx, deleteQuery, recipeID); err != nil {
		return err
	}

	// Add new tags
	for _, tagName := range tagNames {
		tagName = strings.TrimSpace(strings.ToLower(tagName))
		if tagName == "" {
			continue
		}

		// Get or create tag
		tagID, err := getOrCreateTag(ctx, tagName)
		if err != nil {
			return err
		}

		// Link tag to recipe
		insertQuery := `INSERT OR IGNORE INTO recipe_tags (recipe_id, tag_id) VALUES (?, ?)`
		if _, err := db.ExecContext(ctx, insertQuery, recipeID, tagID); err != nil {
			return err
		}
	}

	return nil
}

// getOrCreateTag gets a tag ID by name, creating it if it doesn't exist
func getOrCreateTag(ctx context.Context, tagName string) (string, error) {
	// Try to get existing tag
	var tagID string
	query := `SELECT id FROM tags WHERE name = ?`
	err := db.QueryRowContext(ctx, query, tagName).Scan(&tagID)

	if err == sql.ErrNoRows {
		// Tag doesn't exist, create it
		tagID = uuid.New().String()
		insertQuery := `INSERT INTO tags (id, name) VALUES (?, ?)`
		_, err := db.ExecContext(ctx, insertQuery, tagID, tagName)
		if err != nil {
			return "", err
		}
		return tagID, nil
	}

	if err != nil {
		return "", err
	}

	return tagID, nil
}

// Helper functions for image management

// getRecipeImages returns all images for a given recipe
func getRecipeImages(ctx context.Context, recipeID string) ([]RecipeImage, error) {
	query := `
		SELECT id, recipe_id, image_url, display_order, created_at
		FROM recipe_images
		WHERE recipe_id = ?
		ORDER BY display_order, created_at
	`

	rows, err := db.QueryContext(ctx, query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []RecipeImage
	for rows.Next() {
		var img RecipeImage
		if err := rows.Scan(&img.ID, &img.RecipeID, &img.ImageURL, &img.DisplayOrder, &img.CreatedAt); err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return images, rows.Err()
}

// addRecipeImage adds an image to a recipe
func addRecipeImage(ctx context.Context, recipeID, imageURL string, displayOrder int) error {
	id := uuid.New().String()
	query := `INSERT INTO recipe_images (id, recipe_id, image_url, display_order) VALUES (?, ?, ?, ?)`
	_, err := db.ExecContext(ctx, query, id, recipeID, imageURL, displayOrder)
	return err
}

// deleteRecipeImage deletes an image by ID
func deleteRecipeImage(ctx context.Context, imageID string) error {
	query := `DELETE FROM recipe_images WHERE id = ?`
	_, err := db.ExecContext(ctx, query, imageID)
	return err
}

// User management functions

// DBUser represents a user in the SQLite database
type DBUser struct {
	FirebaseUID string    `json:"firebaseUid"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
	LastLoginAt time.Time `json:"lastLoginAt"`
}

// GetUserByUID retrieves a user from SQLite by Firebase UID
func GetUserByUID(ctx context.Context, firebaseUID string) (*DBUser, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	query := `SELECT firebase_uid, email, display_name, role, created_at, last_login_at FROM users WHERE firebase_uid = ?`

	var user DBUser
	var displayName sql.NullString
	err := db.QueryRowContext(ctx, query, firebaseUID).Scan(
		&user.FirebaseUID, &user.Email, &displayName, &user.Role, &user.CreatedAt, &user.LastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // User not found
	}
	if err != nil {
		return nil, err
	}

	// Handle nullable display_name
	if displayName.Valid {
		user.DisplayName = displayName.String
	} else {
		user.DisplayName = ""
	}

	return &user, nil
}

// CreateUser creates a new user in SQLite
func CreateUser(ctx context.Context, firebaseUID, email, role string) (*DBUser, error) {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	now := time.Now()
	query := `INSERT INTO users (firebase_uid, email, role, created_at, last_login_at) VALUES (?, ?, ?, ?, ?)`

	_, err := db.ExecContext(ctx, query, firebaseUID, email, role, now, now)
	if err != nil {
		return nil, err
	}

	// Upload to Cloud Storage
	if err := uploadDBToGCS(context.Background()); err != nil {
		log.Printf("Failed to upload database to Cloud Storage: %v", err)
	}

	return &DBUser{
		FirebaseUID: firebaseUID,
		Email:       email,
		DisplayName: "",
		Role:        role,
		CreatedAt:   now,
		LastLoginAt: now,
	}, nil
}

// UpdateUserDisplayName updates a user's display name
func UpdateUserDisplayName(ctx context.Context, firebaseUID, displayName string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	query := `UPDATE users SET display_name = ? WHERE firebase_uid = ?`
	_, err := db.ExecContext(ctx, query, displayName, firebaseUID)
	if err != nil {
		return err
	}

	// Upload to Cloud Storage
	if err := uploadDBToGCS(context.Background()); err != nil {
		log.Printf("Failed to upload database to Cloud Storage: %v", err)
	}

	return nil
}

// UpdateUserLastLogin updates a user's last login timestamp
func UpdateUserLastLogin(ctx context.Context, firebaseUID string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	query := `UPDATE users SET last_login_at = ? WHERE firebase_uid = ?`
	_, err := db.ExecContext(ctx, query, time.Now(), firebaseUID)
	if err != nil {
		return err
	}

	// Upload to Cloud Storage
	if err := uploadDBToGCS(context.Background()); err != nil {
		log.Printf("Failed to upload database to Cloud Storage: %v", err)
	}

	return nil
}

// UpdateUserRole updates a user's role (admin only)
func UpdateUserRole(ctx context.Context, firebaseUID, role string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	query := `UPDATE users SET role = ? WHERE firebase_uid = ?`
	_, err := db.ExecContext(ctx, query, role, firebaseUID)
	if err != nil {
		return err
	}

	// Upload to Cloud Storage
	if err := uploadDBToGCS(context.Background()); err != nil {
		log.Printf("Failed to upload database to Cloud Storage: %v", err)
	}

	return nil
}

// Make log management functions

// GetMakeLogsByRecipe returns all make logs for a given recipe
func GetMakeLogsByRecipe(ctx context.Context, recipeID string) ([]MakeLog, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	query := `
		SELECT id, recipe_id, made_at, notes, created_by_user_id, created_at
		FROM make_logs
		WHERE recipe_id = ?
		ORDER BY made_at DESC, created_at DESC
	`

	rows, err := db.QueryContext(ctx, query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []MakeLog
	for rows.Next() {
		var log MakeLog
		var notes sql.NullString
		var createdByUserID sql.NullString

		if err := rows.Scan(&log.ID, &log.RecipeID, &log.MadeAt, &notes, &createdByUserID, &log.CreatedAt); err != nil {
			return nil, err
		}

		if notes.Valid {
			log.Notes = notes.String
		}
		if createdByUserID.Valid {
			uid := createdByUserID.String
			log.CreatedByUserID = &uid
		}

		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// GetMakeCountByRecipe returns the count of make logs for a given recipe
func GetMakeCountByRecipe(ctx context.Context, recipeID string) (int, error) {
	query := `SELECT COUNT(*) FROM make_logs WHERE recipe_id = ?`
	var count int
	err := db.QueryRowContext(ctx, query, recipeID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateMakeLog creates a new make log entry
func CreateMakeLog(ctx context.Context, makeLog *MakeLog) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Generate UUID if not provided
	if makeLog.ID == "" {
		makeLog.ID = uuid.New().String()
	}

	query := `
		INSERT INTO make_logs (id, recipe_id, made_at, notes, created_by_user_id)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := db.ExecContext(ctx, query, makeLog.ID, makeLog.RecipeID, makeLog.MadeAt, makeLog.Notes, makeLog.CreatedByUserID)
	if err != nil {
		return err
	}

	// Upload to Cloud Storage
	if err := uploadDBToGCS(context.Background()); err != nil {
		log.Printf("Failed to upload database to Cloud Storage: %v", err)
	}

	return nil
}

// UpdateMakeLog updates a make log entry
func UpdateMakeLog(ctx context.Context, makeLog *MakeLog) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	query := `
		UPDATE make_logs
		SET made_at = ?, notes = ?
		WHERE id = ?
	`

	result, err := db.ExecContext(ctx, query, makeLog.MadeAt, makeLog.Notes, makeLog.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("make log not found")
	}

	// Upload to Cloud Storage
	if err := uploadDBToGCS(context.Background()); err != nil {
		log.Printf("Failed to upload database to Cloud Storage: %v", err)
	}

	return nil
}

// DeleteMakeLog deletes a make log entry
func DeleteMakeLog(ctx context.Context, logID string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	query := `DELETE FROM make_logs WHERE id = ?`
	result, err := db.ExecContext(ctx, query, logID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("make log not found")
	}

	// Upload to Cloud Storage
	if err := uploadDBToGCS(context.Background()); err != nil {
		log.Printf("Failed to upload database to Cloud Storage: %v", err)
	}

	return nil
}
