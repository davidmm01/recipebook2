package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// YAMLRecipe represents the structure of the YAML files
type YAMLRecipe struct {
	Name        string   `yaml:"name"`
	DateAdded   string   `yaml:"date_added"`
	Source      Source   `yaml:"source"`
	Type        string   `yaml:"type"`
	Descriptors []string `yaml:"descriptors"`
	Cuisine     string   `yaml:"cuisine"`
	Ingredients []string `yaml:"ingredients"`
	Instructions []string `yaml:"instructions"`
	Notes       []string `yaml:"notes"`
	Next        []string `yaml:"next"`
}

type Source struct {
	Name          string `yaml:"name"`
	URL           string `yaml:"url"`
	Type          string `yaml:"type"`
	Modifications string `yaml:"modifications"`
	Submitter     string `yaml:"submitter"`
}

// Recipe represents a recipe for database insertion
type Recipe struct {
	ID              string
	Title           string
	Description     string
	RecipeType      string
	Cuisine         string
	Ingredients     string
	Method          string
	Notes           string
	Sources         string
	Tags            []string
	CreatedByName   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-backend-directory>")
		fmt.Println("Example: go run main.go .")
		os.Exit(1)
	}

	backendPath := os.Args[1]

	// Process cocktails
	cocktailsPath := filepath.Join(backendPath, "data", "cocktails")
	recipesPath := filepath.Join(backendPath, "data", "recipes")

	var allRecipes []Recipe

	// Load cocktails
	cocktails, err := loadRecipesFromDir(cocktailsPath, "drink")
	if err != nil {
		log.Fatalf("Failed to load cocktails: %v", err)
	}
	allRecipes = append(allRecipes, cocktails...)

	// Load recipes
	recipes, err := loadRecipesFromDir(recipesPath, "food")
	if err != nil {
		log.Fatalf("Failed to load recipes: %v", err)
	}
	allRecipes = append(allRecipes, recipes...)

	// Output SQL statements
	ctx := context.Background()
	outputSQL(ctx, allRecipes)

	fmt.Printf("\nSuccessfully processed %d recipes\n", len(allRecipes))
	fmt.Println("\nTo import into your database, run:")
	fmt.Println("  sqlite3 /tmp/recipes.db < import.sql")
}

func loadRecipesFromDir(dirPath string, recipeType string) ([]Recipe, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	var recipes []Recipe

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}

		// Skip template files
		if file.Name() == "template.yaml" {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		recipe, err := loadRecipe(filePath, recipeType)
		if err != nil {
			log.Printf("Warning: failed to load %s: %v", filePath, err)
			continue
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func loadRecipe(filePath string, recipeType string) (Recipe, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Recipe{}, err
	}

	var yamlRecipe YAMLRecipe
	if err := yaml.Unmarshal(data, &yamlRecipe); err != nil {
		return Recipe{}, err
	}

	// Parse the file again to extract comments for ingredients and instructions
	ingredients := parseListWithComments(string(data), "ingredients:")
	instructions := parseListWithComments(string(data), "instructions:")

	// Parse date
	createdAt, err := time.Parse("2006-01-02", yamlRecipe.DateAdded)
	if err != nil {
		// If date parsing fails, use current time
		log.Printf("Warning: failed to parse date %s, using current time", yamlRecipe.DateAdded)
		createdAt = time.Now()
	}

	recipe := Recipe{
		ID:            uuid.New().String(),
		Title:         yamlRecipe.Name,
		Description:   "", // Empty as per instructions
		RecipeType:    recipeType,
		Cuisine:       yamlRecipe.Cuisine,
		Ingredients:   ingredients,
		Method:        instructions,
		Notes:         combineNotes(yamlRecipe.Notes, yamlRecipe.Next),
		Sources:       formatSources(yamlRecipe.Source),
		Tags:          yamlRecipe.Descriptors,
		CreatedByName: yamlRecipe.Source.Submitter,
		CreatedAt:     createdAt,
		UpdatedAt:     time.Now(),
	}

	return recipe, nil
}

// parseListWithComments parses a YAML list field while preserving comments as subheadings
func parseListWithComments(content, fieldName string) string {
	lines := strings.Split(content, "\n")
	var builder strings.Builder
	inSection := false
	baseIndent := -1

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if we're entering the section
		if strings.HasPrefix(trimmed, fieldName) {
			inSection = true
			continue
		}

		if !inSection {
			continue
		}

		// Stop if we hit another field at the same level (not indented)
		if len(line) > 0 && line[0] != ' ' && line[0] != '\t' && line[0] != '-' && line[0] != '#' {
			break
		}

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Get indent level for this line
		indent := len(line) - len(strings.TrimLeft(line, " \t"))

		// If this is the first indented line, set the base indent
		if baseIndent == -1 && indent > 0 {
			baseIndent = indent
		}

		// If we're back to a lower indent level, we've left the section
		if baseIndent > 0 && indent < baseIndent && trimmed != "" {
			break
		}

		// Handle comment lines (subheadings)
		if strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "# ") {
			// This is a YAML comment, skip it
			continue
		}
		if strings.HasPrefix(trimmed, "# ") {
			heading := strings.TrimSpace(strings.TrimPrefix(trimmed, "#"))
			if builder.Len() > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(fmt.Sprintf("### %s\n", heading))
			continue
		}

		// Handle list items
		if strings.HasPrefix(trimmed, "- ") {
			item := strings.TrimPrefix(trimmed, "- ")
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
	}

	return strings.TrimSpace(builder.String())
}

// convertListToMarkdown converts a list of strings to markdown format
// Handles comments (lines starting with #) as subheadings
func convertListToMarkdown(items []string) string {
	if len(items) == 0 {
		return ""
	}

	var builder strings.Builder

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		// Check if it's a comment (subheading)
		if strings.HasPrefix(item, "#") {
			// Convert comment to markdown subheading
			heading := strings.TrimSpace(strings.TrimPrefix(item, "#"))
			if builder.Len() > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(fmt.Sprintf("### %s\n", heading))
		} else {
			// Regular list item
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
	}

	return strings.TrimSpace(builder.String())
}

// combineNotes combines notes and next sections with subheadings
func combineNotes(notes, next []string) string {
	var builder strings.Builder

	if len(notes) > 0 {
		builder.WriteString("### Notes\n")
		for _, note := range notes {
			note = strings.TrimSpace(note)
			if note != "" {
				builder.WriteString(fmt.Sprintf("- %s\n", note))
			}
		}
	}

	if len(next) > 0 {
		if builder.Len() > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString("### Next\n")
		for _, item := range next {
			item = strings.TrimSpace(item)
			if item != "" {
				builder.WriteString(fmt.Sprintf("- %s\n", item))
			}
		}
	}

	return strings.TrimSpace(builder.String())
}

// formatSources formats source information as markdown
func formatSources(source Source) string {
	var parts []string

	if source.Name != "" {
		parts = append(parts, fmt.Sprintf("**Name:** %s", source.Name))
	}
	if source.URL != "" {
		parts = append(parts, fmt.Sprintf("**URL:** %s", source.URL))
	}
	if source.Type != "" {
		parts = append(parts, fmt.Sprintf("**Type:** %s", source.Type))
	}
	if source.Modifications != "" {
		parts = append(parts, fmt.Sprintf("**Modifications:** %s", source.Modifications))
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, "\n")
}

func outputSQL(ctx context.Context, recipes []Recipe) {
	sqlFile, err := os.Create("import.sql")
	if err != nil {
		log.Fatalf("Failed to create import.sql: %v", err)
	}
	defer sqlFile.Close()

	// Write header
	sqlFile.WriteString("-- Auto-generated SQL import script\n")
	sqlFile.WriteString("-- Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")
	sqlFile.WriteString("BEGIN TRANSACTION;\n\n")

	for _, recipe := range recipes {
		// Insert recipe
		sqlFile.WriteString(fmt.Sprintf("-- Recipe: %s\n", recipe.Title))

		// Escape single quotes for SQL
		title := escapeSQLString(recipe.Title)
		description := escapeSQLString(recipe.Description)
		recipeType := escapeSQLString(recipe.RecipeType)
		cuisine := escapeSQLString(recipe.Cuisine)
		ingredients := escapeSQLString(recipe.Ingredients)
		method := escapeSQLString(recipe.Method)
		notes := escapeSQLString(recipe.Notes)
		sources := escapeSQLString(recipe.Sources)
		createdByName := escapeSQLString(recipe.CreatedByName)

		sqlFile.WriteString(fmt.Sprintf(
			"INSERT INTO recipes (id, title, description, recipe_type, cuisine, ingredients, method, notes, sources, created_by_name, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');\n",
			recipe.ID,
			title,
			description,
			recipeType,
			cuisine,
			ingredients,
			method,
			notes,
			sources,
			createdByName,
			recipe.CreatedAt.Format("2006-01-02 15:04:05"),
			recipe.UpdatedAt.Format("2006-01-02 15:04:05"),
		))

		// Insert tags
		for _, tag := range recipe.Tags {
			tag = strings.TrimSpace(strings.ToLower(tag))
			if tag == "" {
				continue
			}
			tagEscaped := escapeSQLString(tag)

			// Insert tag if it doesn't exist
			sqlFile.WriteString(fmt.Sprintf("INSERT OR IGNORE INTO tags (name) VALUES ('%s');\n", tagEscaped))

			// Link tag to recipe
			sqlFile.WriteString(fmt.Sprintf(
				"INSERT INTO recipe_tags (recipe_id, tag_id) SELECT '%s', id FROM tags WHERE name = '%s';\n",
				recipe.ID, tagEscaped,
			))
		}

		sqlFile.WriteString("\n")
	}

	sqlFile.WriteString("COMMIT;\n")

	fmt.Println("Generated import.sql file")
}

func escapeSQLString(s string) string {
	// Replace single quotes with two single quotes for SQL escaping
	return strings.ReplaceAll(s, "'", "''")
}
