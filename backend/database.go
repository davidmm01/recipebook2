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
	ID         int64     `json:"id"`
	Filename   string    `json:"filename"`
	IconURL    string    `json:"iconUrl"`
	UploadedAt time.Time `json:"uploadedAt"`
}

// Recipe represents a recipe with markdown fields
type Recipe struct {
	ID          string        `json:"id"` // UUID
	Title       string        `json:"title"`
	Description string        `json:"description"` // Brief description
	RecipeType  string        `json:"type"`        // "food", "cocktail", etc.
	Cuisine     string        `json:"cuisine"`     // "italian", "japanese", "mexican", etc.
	Ingredients string        `json:"ingredients"` // markdown
	Method      string        `json:"method"`      // markdown
	Notes       string        `json:"notes"`       // markdown
	Sources     string        `json:"sources"`     // markdown
	IconID      *int64        `json:"iconId"`      // Nullable icon ID
	Icon        *Icon         `json:"icon"`        // Icon details (loaded separately)
	Tags        []string      `json:"tags"`        // Array of tag names
	Images      []RecipeImage `json:"images"`      // Array of image URLs
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

// RecipeImage represents an image associated with a recipe
type RecipeImage struct {
	ID           int64     `json:"id"`
	RecipeID     string    `json:"recipeId"`
	ImageURL     string    `json:"imageUrl"`
	DisplayOrder int       `json:"displayOrder"`
	CreatedAt    time.Time `json:"createdAt"`
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
	-- Icons table (shared collection of recipe icons)
	CREATE TABLE IF NOT EXISTS icons (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
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
		icon_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (icon_id) REFERENCES icons(id)
	);

	CREATE INDEX IF NOT EXISTS idx_recipes_type ON recipes(recipe_type);
	CREATE INDEX IF NOT EXISTS idx_recipes_cuisine ON recipes(cuisine);
	CREATE INDEX IF NOT EXISTS idx_recipes_icon ON recipes(icon_id);

	-- Tags table
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
	);

	-- Many-to-many relationship between recipes and tags
	CREATE TABLE IF NOT EXISTS recipe_tags (
		recipe_id TEXT NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (recipe_id, tag_id),
		FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_recipe_tags_recipe ON recipe_tags(recipe_id);
	CREATE INDEX IF NOT EXISTS idx_recipe_tags_tag ON recipe_tags(tag_id);

	-- Recipe images table
	CREATE TABLE IF NOT EXISTS recipe_images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		recipe_id TEXT NOT NULL,
		image_url TEXT NOT NULL,
		display_order INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_recipe_images_recipe ON recipe_images(recipe_id);

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
		SELECT id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, icon_id, created_at, updated_at
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
		err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.Sources, &r.IconID, &r.CreatedAt, &r.UpdatedAt)
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

		recipes = append(recipes, r)
	}

	return recipes, rows.Err()
}

// GetRecipeByID returns a single recipe by ID
func GetRecipeByID(ctx context.Context, recipeID string) (*Recipe, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	query := `
		SELECT id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, icon_id, created_at, updated_at
		FROM recipes
		WHERE id = ?
	`

	var r Recipe
	err := db.QueryRowContext(ctx, query, recipeID).Scan(
		&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.Sources, &r.IconID, &r.CreatedAt, &r.UpdatedAt,
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

	return &r, nil
}

// SearchRecipes performs full-text search across recipes
func SearchRecipes(ctx context.Context, query string) ([]Recipe, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	sqlQuery := `
		SELECT r.id, r.title, r.description, r.recipe_type, r.cuisine, r.ingredients, r.method, r.notes, r.sources, r.icon_id, r.created_at, r.updated_at
		FROM recipes r
		JOIN recipes_fts ON r.id = recipes_fts.recipe_id
		WHERE recipes_fts MATCH ?
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
		err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.Sources, &r.IconID, &r.CreatedAt, &r.UpdatedAt)
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

		recipes = append(recipes, r)
	}

	return recipes, rows.Err()
}

// FilterRecipes performs filtering based on search text, tags, and cuisine
// If searchQuery is provided, it uses FTS5 for text search
// If tags are provided, filters recipes that have ALL specified tags
// If cuisine is provided, filters by exact cuisine match
func FilterRecipes(ctx context.Context, searchQuery string, tags []string, cuisine string) ([]Recipe, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	var queryBuilder strings.Builder
	var args []interface{}

	// Base query - start with recipes table
	if searchQuery != "" {
		// Use FTS5 for text search
		queryBuilder.WriteString(`
			SELECT DISTINCT r.id, r.title, r.description, r.recipe_type, r.cuisine, r.ingredients, r.method, r.notes, r.sources, r.icon_id, r.created_at, r.updated_at
			FROM recipes r
			JOIN recipes_fts ON r.id = recipes_fts.recipe_id
			WHERE recipes_fts MATCH ?
		`)
		args = append(args, searchQuery)
	} else {
		// No text search, just filter
		queryBuilder.WriteString(`
			SELECT DISTINCT r.id, r.title, r.description, r.recipe_type, r.cuisine, r.ingredients, r.method, r.notes, r.sources, r.icon_id, r.created_at, r.updated_at
			FROM recipes r
			WHERE 1=1
		`)
	}

	// Add cuisine filter
	if cuisine != "" {
		if searchQuery != "" {
			queryBuilder.WriteString(` AND r.cuisine = ?`)
		} else {
			queryBuilder.WriteString(` AND r.cuisine = ?`)
		}
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
	if searchQuery != "" {
		queryBuilder.WriteString(` ORDER BY rank`)
	} else {
		queryBuilder.WriteString(` ORDER BY r.updated_at DESC`)
	}

	rows, err := db.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.Sources, &r.IconID, &r.CreatedAt, &r.UpdatedAt)
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
		INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, icon_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.ExecContext(ctx, query,
		recipe.ID, recipe.Title, recipe.Description, recipe.RecipeType, recipe.Cuisine,
		recipe.Ingredients, recipe.Method, recipe.Notes, recipe.Sources, recipe.IconID, recipe.CreatedAt, recipe.UpdatedAt,
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

	// Upload to Cloud Storage (async)
	go func() {
		if err := uploadDBToGCS(context.Background()); err != nil {
			log.Printf("Failed to upload database to Cloud Storage: %v", err)
		}
	}()

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

	// Upload to Cloud Storage (async)
	go func() {
		if err := uploadDBToGCS(context.Background()); err != nil {
			log.Printf("Failed to upload database to Cloud Storage: %v", err)
		}
	}()

	return nil
}

// Helper functions for icon management

// getIconByID returns an icon by ID
func getIconByID(ctx context.Context, iconID int64) (*Icon, error) {
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
func createIcon(ctx context.Context, filename, iconURL string) (int64, error) {
	query := `INSERT INTO icons (filename, icon_url) VALUES (?, ?)`
	result, err := db.ExecContext(ctx, query, filename, iconURL)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAllTags returns all unique tags in the database
func GetAllTags(ctx context.Context) ([]string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	query := `SELECT name FROM tags ORDER BY name`
	rows, err := db.QueryContext(ctx, query)
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

// GetAllCuisines returns all unique cuisines in the database
func GetAllCuisines(ctx context.Context) ([]string, error) {
	dbMutex.RLock()
	defer dbMutex.RUnlock()

	query := `SELECT DISTINCT cuisine FROM recipes WHERE cuisine IS NOT NULL AND cuisine != '' ORDER BY cuisine`
	rows, err := db.QueryContext(ctx, query)
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
func getOrCreateTag(ctx context.Context, tagName string) (int64, error) {
	// Try to get existing tag
	var tagID int64
	query := `SELECT id FROM tags WHERE name = ?`
	err := db.QueryRowContext(ctx, query, tagName).Scan(&tagID)

	if err == sql.ErrNoRows {
		// Tag doesn't exist, create it
		insertQuery := `INSERT INTO tags (name) VALUES (?)`
		result, err := db.ExecContext(ctx, insertQuery, tagName)
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()
	}

	if err != nil {
		return 0, err
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
	query := `INSERT INTO recipe_images (recipe_id, image_url, display_order) VALUES (?, ?, ?)`
	_, err := db.ExecContext(ctx, query, recipeID, imageURL, displayOrder)
	return err
}

// deleteRecipeImage deletes an image by ID
func deleteRecipeImage(ctx context.Context, imageID int64) error {
	query := `DELETE FROM recipe_images WHERE id = ?`
	_, err := db.ExecContext(ctx, query, imageID)
	return err
}
