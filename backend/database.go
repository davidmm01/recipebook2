package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/storage"
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

// Recipe represents a recipe with markdown fields
type Recipe struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	RecipeType  string    `json:"type"` // "food", "cocktail", etc.
	Ingredients string    `json:"ingredients"` // markdown
	Method      string    `json:"method"` // markdown
	Notes       string    `json:"notes"` // markdown
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
	CREATE TABLE IF NOT EXISTS recipes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		recipe_type TEXT,
		ingredients TEXT,
		method TEXT,
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_recipes_type ON recipes(recipe_type);

	-- Full-text search table
	CREATE VIRTUAL TABLE IF NOT EXISTS recipes_fts USING fts5(
		title, ingredients, method, notes,
		content=recipes,
		content_rowid=id
	);

	-- Triggers to keep FTS in sync
	CREATE TRIGGER IF NOT EXISTS recipes_fts_insert AFTER INSERT ON recipes BEGIN
		INSERT INTO recipes_fts(rowid, title, ingredients, method, notes)
		VALUES (new.id, new.title, new.ingredients, new.method, new.notes);
	END;

	CREATE TRIGGER IF NOT EXISTS recipes_fts_update AFTER UPDATE ON recipes BEGIN
		UPDATE recipes_fts
		SET title=new.title, ingredients=new.ingredients,
			method=new.method, notes=new.notes
		WHERE rowid=new.id;
	END;

	CREATE TRIGGER IF NOT EXISTS recipes_fts_delete AFTER DELETE ON recipes BEGIN
		DELETE FROM recipes_fts WHERE rowid=old.id;
	END;
	`

	if _, err := database.Exec(schema); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	log.Println("Created new database with schema")
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

// uploadDBToGCS uploads the SQLite database to Cloud Storage
func uploadDBToGCS(ctx context.Context) error {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

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
		SELECT id, title, recipe_type, ingredients, method, notes, created_at, updated_at
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
		err := rows.Scan(&r.ID, &r.Title, &r.RecipeType, &r.Ingredients, &r.Method, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, rows.Err()
}

// GetRecipeByID returns a single recipe by ID
func GetRecipeByID(ctx context.Context, recipeID int64) (*Recipe, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	query := `
		SELECT id, title, recipe_type, ingredients, method, notes, created_at, updated_at
		FROM recipes
		WHERE id = ?
	`

	var r Recipe
	err := db.QueryRowContext(ctx, query, recipeID).Scan(
		&r.ID, &r.Title, &r.RecipeType, &r.Ingredients, &r.Method, &r.Notes, &r.CreatedAt, &r.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// SearchRecipes performs full-text search across recipes
func SearchRecipes(ctx context.Context, query string) ([]Recipe, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	sqlQuery := `
		SELECT r.id, r.title, r.recipe_type, r.ingredients, r.method, r.notes, r.created_at, r.updated_at
		FROM recipes r
		JOIN recipes_fts fts ON r.id = fts.rowid
		WHERE fts MATCH ?
		ORDER BY rank
	`

	rows, err := db.QueryContext(ctx, sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		err := rows.Scan(&r.ID, &r.Title, &r.RecipeType, &r.Ingredients, &r.Method, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, r)
	}

	return recipes, rows.Err()
}

// CreateRecipe inserts a new recipe and syncs to Cloud Storage
func CreateRecipe(ctx context.Context, recipe *Recipe) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	recipe.CreatedAt = time.Now()
	recipe.UpdatedAt = time.Now()

	query := `
		INSERT INTO recipes (title, recipe_type, ingredients, method, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.ExecContext(ctx, query,
		recipe.Title, recipe.RecipeType, recipe.Ingredients, recipe.Method, recipe.Notes, recipe.CreatedAt, recipe.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	recipe.ID = id

	// Upload to Cloud Storage (async to not block response)
	go func() {
		if err := uploadDBToGCS(context.Background()); err != nil {
			log.Printf("Failed to upload database to Cloud Storage: %v", err)
		}
	}()

	return nil
}

// UpdateRecipe updates an existing recipe and syncs to Cloud Storage
func UpdateRecipe(ctx context.Context, recipe *Recipe) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	recipe.UpdatedAt = time.Now()

	query := `
		UPDATE recipes
		SET title = ?, recipe_type = ?, ingredients = ?, method = ?, notes = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := db.ExecContext(ctx, query,
		recipe.Title, recipe.RecipeType, recipe.Ingredients, recipe.Method, recipe.Notes, recipe.UpdatedAt,
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

	// Upload to Cloud Storage (async)
	go func() {
		if err := uploadDBToGCS(context.Background()); err != nil {
			log.Printf("Failed to upload database to Cloud Storage: %v", err)
		}
	}()

	return nil
}

// DeleteRecipe deletes a recipe and syncs to Cloud Storage
func DeleteRecipe(ctx context.Context, recipeID int64) error {
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

	// Upload to Cloud Storage (async)
	go func() {
		if err := uploadDBToGCS(context.Background()); err != nil {
			log.Printf("Failed to upload database to Cloud Storage: %v", err)
		}
	}()

	return nil
}
