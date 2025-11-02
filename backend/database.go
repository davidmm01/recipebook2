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

// Recipe represents a recipe with markdown fields
type Recipe struct {
	ID          string        `json:"id"`          // UUID
	Title       string        `json:"title"`
	Description string        `json:"description"` // Brief description
	RecipeType  string        `json:"type"`        // "food", "cocktail", etc.
	Cuisine     string        `json:"cuisine"`     // "italian", "japanese", "mexican", etc.
	Ingredients string        `json:"ingredients"` // markdown
	Method      string        `json:"method"`      // markdown
	Notes       string        `json:"notes"`       // markdown
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
	CREATE TABLE IF NOT EXISTS recipes (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		recipe_type TEXT,
		cuisine TEXT,
		ingredients TEXT,
		method TEXT,
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_recipes_type ON recipes(recipe_type);
	CREATE INDEX IF NOT EXISTS idx_recipes_cuisine ON recipes(cuisine);

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

	-- Full-text search table (now includes description and cuisine)
	-- Note: We include recipe_id as an unindexed column to enable joining back to recipes table
	CREATE VIRTUAL TABLE IF NOT EXISTS recipes_fts USING fts5(
		recipe_id UNINDEXED,
		title, description, cuisine, ingredients, method, notes
	);

	-- Triggers to keep FTS in sync
	CREATE TRIGGER IF NOT EXISTS recipes_fts_insert AFTER INSERT ON recipes BEGIN
		INSERT INTO recipes_fts(recipe_id, title, description, cuisine, ingredients, method, notes)
		VALUES (new.id, new.title, new.description, new.cuisine, new.ingredients, new.method, new.notes);
	END;

	CREATE TRIGGER IF NOT EXISTS recipes_fts_update AFTER UPDATE ON recipes BEGIN
		UPDATE recipes_fts
		SET title=new.title, description=new.description, cuisine=new.cuisine,
			ingredients=new.ingredients, method=new.method, notes=new.notes
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
		SELECT id, title, description, recipe_type, cuisine, ingredients, method, notes, created_at, updated_at
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
		err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
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
		SELECT id, title, description, recipe_type, cuisine, ingredients, method, notes, created_at, updated_at
		FROM recipes
		WHERE id = ?
	`

	var r Recipe
	err := db.QueryRowContext(ctx, query, recipeID).Scan(
		&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.CreatedAt, &r.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
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
		SELECT r.id, r.title, r.description, r.recipe_type, r.cuisine, r.ingredients, r.method, r.notes, r.created_at, r.updated_at
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
		err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.RecipeType, &r.Cuisine, &r.Ingredients, &r.Method, &r.Notes, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
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
		INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.ExecContext(ctx, query,
		recipe.ID, recipe.Title, recipe.Description, recipe.RecipeType, recipe.Cuisine,
		recipe.Ingredients, recipe.Method, recipe.Notes, recipe.CreatedAt, recipe.UpdatedAt,
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
		    ingredients = ?, method = ?, notes = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := db.ExecContext(ctx, query,
		recipe.Title, recipe.Description, recipe.RecipeType, recipe.Cuisine,
		recipe.Ingredients, recipe.Method, recipe.Notes, recipe.UpdatedAt,
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
