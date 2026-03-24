package main

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	var err error
	// Use a shared in-memory database so all connections within the pool
	// see the same data. Plain ":memory:" gives each connection its own
	// empty database, which breaks nested queries (e.g. GetRecipes calls
	// getRecipeTags on a second connection that has no schema).
	db, err = sql.Open("sqlite3", "file:testrecipebook?mode=memory&cache=shared")
	if err != nil {
		panic("failed to open in-memory database: " + err.Error())
	}
	defer db.Close()

	// Suppress GCS uploads during tests
	dbSyncDisabled = true

	// Enable foreign key enforcement so ON DELETE CASCADE works
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		panic("failed to enable foreign keys: " + err.Error())
	}

	schemaBytes, err := os.ReadFile("schema.sql")
	if err != nil {
		panic("failed to read schema.sql: " + err.Error())
	}
	if _, err := db.Exec(string(schemaBytes)); err != nil {
		panic("failed to create schema: " + err.Error())
	}

	os.Exit(m.Run())
}

// clearTables removes all rows from all tables before each test.
// recipes_fts is cleared automatically via delete trigger on recipes.
func clearTables(t *testing.T) {
	t.Helper()
	tables := []string{"make_logs", "recipe_tags", "recipe_images", "recipes", "tags", "icons", "users"}
	for _, table := range tables {
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			t.Fatalf("failed to clear table %s: %v", table, err)
		}
	}
}

func makeTestUser(uid, email, role string) *DBUser {
	user, err := CreateUser(context.Background(), uid, email, role)
	if err != nil {
		panic("makeTestUser failed: " + err.Error())
	}
	return user
}

func makeTestRecipe(title string) *Recipe {
	return &Recipe{
		Title:       title,
		Description: "A test recipe",
		RecipeType:  "food",
		Cuisine:     "italian",
		Ingredients: "1 cup flour\n2 eggs",
		Method:      "Mix and cook",
	}
}

// --- Recipe CRUD ---

func TestCreateRecipe_GeneratesUUID(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("UUID Test Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}
	if r.ID == "" {
		t.Error("expected UUID to be generated, got empty string")
	}
}

func TestCreateRecipe_WithTags(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Tagged Recipe")
	r.Tags = []string{"pasta", "vegetarian"}
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	got, err := GetRecipeByID(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	if got == nil {
		t.Fatal("expected recipe, got nil")
	}
	if len(got.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d: %v", len(got.Tags), got.Tags)
	}
}

func TestCreateRecipe_TagsNormalizedToLowercase(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Lowercase Tag Test")
	r.Tags = []string{"Italian", "VEGAN"}
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	got, err := GetRecipeByID(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	for _, tag := range got.Tags {
		if tag != strings.ToLower(tag) {
			t.Errorf("expected tag to be lowercase, got '%s'", tag)
		}
	}
}

func TestGetRecipes_Empty(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	recipes, err := GetRecipes(ctx)
	if err != nil {
		t.Fatalf("GetRecipes failed: %v", err)
	}
	if len(recipes) != 0 {
		t.Errorf("expected 0 recipes, got %d", len(recipes))
	}
}

func TestGetRecipes_Multiple(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r1 := makeTestRecipe("Recipe Alpha")
	r2 := makeTestRecipe("Recipe Beta")
	if err := CreateRecipe(ctx, r1); err != nil {
		t.Fatalf("CreateRecipe r1 failed: %v", err)
	}
	if err := CreateRecipe(ctx, r2); err != nil {
		t.Fatalf("CreateRecipe r2 failed: %v", err)
	}

	recipes, err := GetRecipes(ctx)
	if err != nil {
		t.Fatalf("GetRecipes failed: %v", err)
	}
	if len(recipes) != 2 {
		t.Errorf("expected 2 recipes, got %d", len(recipes))
	}
}

func TestGetRecipeByID_Found(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Find Me")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	got, err := GetRecipeByID(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	if got == nil {
		t.Fatal("expected recipe, got nil")
	}
	if got.Title != "Find Me" {
		t.Errorf("expected title 'Find Me', got '%s'", got.Title)
	}
	if got.Cuisine != "italian" {
		t.Errorf("expected cuisine 'italian', got '%s'", got.Cuisine)
	}
}

func TestGetRecipeByID_NotFound(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	got, err := GetRecipeByID(ctx, "nonexistent-id")
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for missing recipe, got %+v", got)
	}
}

func TestUpdateRecipe_Success(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Original Title")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	r.Title = "Updated Title"
	r.Cuisine = "french"
	if err := UpdateRecipe(ctx, r); err != nil {
		t.Fatalf("UpdateRecipe failed: %v", err)
	}

	got, err := GetRecipeByID(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	if got.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", got.Title)
	}
	if got.Cuisine != "french" {
		t.Errorf("expected cuisine 'french', got '%s'", got.Cuisine)
	}
}

func TestUpdateRecipe_UpdatesTags(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Tag Update Test")
	r.Tags = []string{"old-tag"}
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	r.Tags = []string{"new-tag-1", "new-tag-2"}
	if err := UpdateRecipe(ctx, r); err != nil {
		t.Fatalf("UpdateRecipe failed: %v", err)
	}

	got, err := GetRecipeByID(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	if len(got.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d: %v", len(got.Tags), got.Tags)
	}
	for _, tag := range got.Tags {
		if tag == "old-tag" {
			t.Error("old-tag should have been removed after update")
		}
	}
}

func TestUpdateRecipe_NotFound(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := &Recipe{ID: "nonexistent-id", Title: "Ghost"}
	if err := UpdateRecipe(ctx, r); err == nil {
		t.Error("expected error for nonexistent recipe, got nil")
	}
}

func TestDeleteRecipe_Success(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Delete Me")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	if err := DeleteRecipe(ctx, r.ID); err != nil {
		t.Fatalf("DeleteRecipe failed: %v", err)
	}

	got, err := GetRecipeByID(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil after delete, got %+v", got)
	}
}

func TestDeleteRecipe_CascadesTags(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Cascade Tags Test")
	r.Tags = []string{"cascade-tag"}
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	if err := DeleteRecipe(ctx, r.ID); err != nil {
		t.Fatalf("DeleteRecipe failed: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM recipe_tags WHERE recipe_id = ?", r.ID).Scan(&count); err != nil {
		t.Fatalf("query recipe_tags failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 recipe_tags rows after delete, got %d", count)
	}
}

func TestDeleteRecipe_NotFound(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	if err := DeleteRecipe(ctx, "nonexistent-id"); err == nil {
		t.Error("expected error for nonexistent recipe, got nil")
	}
}

// --- Filter tests ---

func TestFilterRecipes_ByCuisine(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	italian := makeTestRecipe("Italian Pasta")
	italian.Cuisine = "italian"
	french := makeTestRecipe("French Onion Soup")
	french.Cuisine = "french"

	if err := CreateRecipe(ctx, italian); err != nil {
		t.Fatalf("CreateRecipe italian failed: %v", err)
	}
	if err := CreateRecipe(ctx, french); err != nil {
		t.Fatalf("CreateRecipe french failed: %v", err)
	}

	results, err := FilterRecipes(ctx, "", nil, "italian", "", "")
	if err != nil {
		t.Fatalf("FilterRecipes failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if len(results) > 0 && results[0].Cuisine != "italian" {
		t.Errorf("expected italian cuisine, got '%s'", results[0].Cuisine)
	}
}

func TestFilterRecipes_ByTag(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r1 := makeTestRecipe("Vegetarian Dish")
	r1.Tags = []string{"vegetarian"}
	r2 := makeTestRecipe("Meat Dish")
	r2.Tags = []string{"meat"}

	if err := CreateRecipe(ctx, r1); err != nil {
		t.Fatalf("CreateRecipe r1 failed: %v", err)
	}
	if err := CreateRecipe(ctx, r2); err != nil {
		t.Fatalf("CreateRecipe r2 failed: %v", err)
	}

	results, err := FilterRecipes(ctx, "", []string{"vegetarian"}, "", "", "")
	if err != nil {
		t.Fatalf("FilterRecipes failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestFilterRecipes_ByRecipeType(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	food := makeTestRecipe("A Food")
	food.RecipeType = "food"
	cocktail := makeTestRecipe("A Cocktail")
	cocktail.RecipeType = "cocktail"

	if err := CreateRecipe(ctx, food); err != nil {
		t.Fatalf("CreateRecipe food failed: %v", err)
	}
	if err := CreateRecipe(ctx, cocktail); err != nil {
		t.Fatalf("CreateRecipe cocktail failed: %v", err)
	}

	results, err := FilterRecipes(ctx, "", nil, "", "cocktail", "")
	if err != nil {
		t.Fatalf("FilterRecipes failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if len(results) > 0 && results[0].RecipeType != "cocktail" {
		t.Errorf("expected cocktail type, got '%s'", results[0].RecipeType)
	}
}

func TestFilterRecipes_Combined(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r1 := makeTestRecipe("Italian Vegan Pasta")
	r1.Cuisine = "italian"
	r1.Tags = []string{"vegan"}

	r2 := makeTestRecipe("Italian Carbonara")
	r2.Cuisine = "italian"
	r2.Tags = []string{"meat"}

	r3 := makeTestRecipe("French Vegan Soup")
	r3.Cuisine = "french"
	r3.Tags = []string{"vegan"}

	for _, r := range []*Recipe{r1, r2, r3} {
		if err := CreateRecipe(ctx, r); err != nil {
			t.Fatalf("CreateRecipe failed: %v", err)
		}
	}

	results, err := FilterRecipes(ctx, "", []string{"vegan"}, "italian", "", "")
	if err != nil {
		t.Fatalf("FilterRecipes failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if len(results) > 0 && results[0].Title != "Italian Vegan Pasta" {
		t.Errorf("expected 'Italian Vegan Pasta', got '%s'", results[0].Title)
	}
}

// --- Search tests ---

func TestPrepareFTS5Query(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"bour", "bour*"},
		{"lime juice", "lime* juice*"},
		{"OR", "OR"},
		{"AND", "AND"},
		{"NOT", "NOT"},
		{"bour*", "bour*"},
		{"", ""},
	}
	for _, tt := range tests {
		got := prepareFTS5Query(tt.input)
		if got != tt.expected {
			t.Errorf("prepareFTS5Query(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestSearchRecipes_BasicMatch(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Beef Bourguignon")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	results, err := SearchRecipes(ctx, "bour")
	if err != nil {
		t.Fatalf("SearchRecipes failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result for prefix 'bour', got %d", len(results))
	}
}

func TestSearchRecipes_NoMatch(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Chicken Tikka Masala")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	results, err := SearchRecipes(ctx, "bourguignon")
	if err != nil {
		t.Fatalf("SearchRecipes failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchRecipes_SearchesIngredients(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Mystery Dish")
	r.Ingredients = "3 cups lemongrass\n2 tsp galangal"
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	results, err := SearchRecipes(ctx, "lemongrass")
	if err != nil {
		t.Fatalf("SearchRecipes failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result searching by ingredient, got %d", len(results))
	}
}

// --- Tags and cuisines ---

func TestGetAllTagsWithCounts_Empty(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	tags, err := GetAllTagsWithCounts(ctx, "")
	if err != nil {
		t.Fatalf("GetAllTagsWithCounts failed: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 tags, got %d", len(tags))
	}
}

func TestGetAllTagsWithCounts_WithData(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r1 := makeTestRecipe("Recipe 1")
	r1.Tags = []string{"pasta", "italian"}
	r2 := makeTestRecipe("Recipe 2")
	r2.Tags = []string{"pasta"}

	if err := CreateRecipe(ctx, r1); err != nil {
		t.Fatalf("CreateRecipe r1 failed: %v", err)
	}
	if err := CreateRecipe(ctx, r2); err != nil {
		t.Fatalf("CreateRecipe r2 failed: %v", err)
	}

	tags, err := GetAllTagsWithCounts(ctx, "")
	if err != nil {
		t.Fatalf("GetAllTagsWithCounts failed: %v", err)
	}

	tagMap := make(map[string]int)
	for _, tag := range tags {
		tagMap[tag.Name] = tag.Count
	}

	if tagMap["pasta"] != 2 {
		t.Errorf("expected pasta count=2, got %d", tagMap["pasta"])
	}
	if tagMap["italian"] != 1 {
		t.Errorf("expected italian count=1, got %d", tagMap["italian"])
	}
}

func TestGetAllCuisines_Empty(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	cuisines, err := GetAllCuisines(ctx, "")
	if err != nil {
		t.Fatalf("GetAllCuisines failed: %v", err)
	}
	if len(cuisines) != 0 {
		t.Errorf("expected 0 cuisines, got %d", len(cuisines))
	}
}

func TestGetAllCuisines_ReturnsDistinct(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r1 := makeTestRecipe("Pasta 1")
	r1.Cuisine = "italian"
	r2 := makeTestRecipe("Pasta 2")
	r2.Cuisine = "italian"
	r3 := makeTestRecipe("Croissant")
	r3.Cuisine = "french"

	for _, r := range []*Recipe{r1, r2, r3} {
		if err := CreateRecipe(ctx, r); err != nil {
			t.Fatalf("CreateRecipe failed: %v", err)
		}
	}

	cuisines, err := GetAllCuisines(ctx, "")
	if err != nil {
		t.Fatalf("GetAllCuisines failed: %v", err)
	}
	if len(cuisines) != 2 {
		t.Errorf("expected 2 distinct cuisines, got %d: %v", len(cuisines), cuisines)
	}
}

// --- Make log tests ---

func TestCreateMakeLog_Success(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Logged Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	makeLog := &MakeLog{
		RecipeID: r.ID,
		MadeAt:   "2024-01-15",
		Notes:    "Turned out great",
	}
	if err := CreateMakeLog(ctx, makeLog); err != nil {
		t.Fatalf("CreateMakeLog failed: %v", err)
	}
	if makeLog.ID == "" {
		t.Error("expected make log ID to be generated")
	}
}

func TestGetMakeLogsByRecipe_Empty(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("No Logs Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	logs, err := GetMakeLogsByRecipe(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetMakeLogsByRecipe failed: %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("expected 0 logs, got %d", len(logs))
	}
}

func TestGetMakeLogsByRecipe_Multiple(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Multi-Log Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	for _, date := range []string{"2024-01-15", "2024-02-20", "2024-03-01"} {
		if err := CreateMakeLog(ctx, &MakeLog{RecipeID: r.ID, MadeAt: date}); err != nil {
			t.Fatalf("CreateMakeLog failed: %v", err)
		}
	}

	logs, err := GetMakeLogsByRecipe(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetMakeLogsByRecipe failed: %v", err)
	}
	if len(logs) != 3 {
		t.Errorf("expected 3 logs, got %d", len(logs))
	}
}

func TestMakeCountReflectedInRecipe(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Count Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}
	if err := CreateMakeLog(ctx, &MakeLog{RecipeID: r.ID, MadeAt: "2024-01-15"}); err != nil {
		t.Fatalf("CreateMakeLog failed: %v", err)
	}

	got, err := GetRecipeByID(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	if got.MakeCount != 1 {
		t.Errorf("expected MakeCount=1, got %d", got.MakeCount)
	}
}

func TestUpdateMakeLog_Success(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Update Log Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	makeLog := &MakeLog{RecipeID: r.ID, MadeAt: "2024-01-15", Notes: "original note"}
	if err := CreateMakeLog(ctx, makeLog); err != nil {
		t.Fatalf("CreateMakeLog failed: %v", err)
	}

	makeLog.MadeAt = "2024-06-01"
	makeLog.Notes = "updated note"
	if err := UpdateMakeLog(ctx, makeLog); err != nil {
		t.Fatalf("UpdateMakeLog failed: %v", err)
	}

	logs, err := GetMakeLogsByRecipe(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetMakeLogsByRecipe failed: %v", err)
	}
	if len(logs) == 0 {
		t.Fatal("expected logs, got none")
	}
	// go-sqlite3 may return DATE values as "YYYY-MM-DD" or "YYYY-MM-DDT00:00:00Z"
	if !strings.HasPrefix(logs[0].MadeAt, "2024-06-01") {
		t.Errorf("expected date starting with '2024-06-01', got '%s'", logs[0].MadeAt)
	}
	if logs[0].Notes != "updated note" {
		t.Errorf("expected 'updated note', got '%s'", logs[0].Notes)
	}
}

func TestDeleteMakeLog_Success(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Delete Log Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	makeLog := &MakeLog{RecipeID: r.ID, MadeAt: "2024-01-15"}
	if err := CreateMakeLog(ctx, makeLog); err != nil {
		t.Fatalf("CreateMakeLog failed: %v", err)
	}

	if err := DeleteMakeLog(ctx, makeLog.ID); err != nil {
		t.Fatalf("DeleteMakeLog failed: %v", err)
	}

	logs, err := GetMakeLogsByRecipe(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetMakeLogsByRecipe failed: %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("expected 0 logs after delete, got %d", len(logs))
	}
}

// --- User tests ---

func TestCreateUser_Success(t *testing.T) {
	clearTables(t)

	user := makeTestUser("uid-1", "alice@example.com", "viewer")
	if user.FirebaseUID != "uid-1" {
		t.Errorf("expected UID 'uid-1', got '%s'", user.FirebaseUID)
	}
	if user.Role != "viewer" {
		t.Errorf("expected role 'viewer', got '%s'", user.Role)
	}
}

func TestGetUserByUID_Found(t *testing.T) {
	clearTables(t)
	makeTestUser("uid-1", "alice@example.com", "editor")

	user, err := GetUserByUID(context.Background(), "uid-1")
	if err != nil {
		t.Fatalf("GetUserByUID failed: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected email 'alice@example.com', got '%s'", user.Email)
	}
	if user.Role != "editor" {
		t.Errorf("expected role 'editor', got '%s'", user.Role)
	}
}

func TestGetUserByUID_NotFound(t *testing.T) {
	clearTables(t)

	user, err := GetUserByUID(context.Background(), "nonexistent-uid")
	if err != nil {
		t.Fatalf("GetUserByUID failed: %v", err)
	}
	if user != nil {
		t.Errorf("expected nil for missing user, got %+v", user)
	}
}

func TestGetAllUsers_Empty(t *testing.T) {
	clearTables(t)

	users, err := GetAllUsers(context.Background())
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}

func TestGetAllUsers_ReturnsAll(t *testing.T) {
	clearTables(t)
	makeTestUser("uid-1", "alice@example.com", "viewer")
	makeTestUser("uid-2", "bob@example.com", "editor")
	makeTestUser("uid-3", "carol@example.com", "admin")

	users, err := GetAllUsers(context.Background())
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestGetAllUsers_OrderedByEmail(t *testing.T) {
	clearTables(t)
	makeTestUser("uid-1", "zara@example.com", "viewer")
	makeTestUser("uid-2", "alice@example.com", "viewer")
	makeTestUser("uid-3", "mike@example.com", "viewer")

	users, err := GetAllUsers(context.Background())
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}
	if users[0].Email != "alice@example.com" {
		t.Errorf("expected first user to be alice@example.com, got '%s'", users[0].Email)
	}
	if users[2].Email != "zara@example.com" {
		t.Errorf("expected last user to be zara@example.com, got '%s'", users[2].Email)
	}
}

func TestUpdateUserRole_Success(t *testing.T) {
	clearTables(t)
	makeTestUser("uid-1", "alice@example.com", "viewer")

	if err := UpdateUserRole(context.Background(), "uid-1", "admin"); err != nil {
		t.Fatalf("UpdateUserRole failed: %v", err)
	}

	user, err := GetUserByUID(context.Background(), "uid-1")
	if err != nil {
		t.Fatalf("GetUserByUID failed: %v", err)
	}
	if user.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", user.Role)
	}
}

func TestUpdateUserDisplayName_Success(t *testing.T) {
	clearTables(t)
	makeTestUser("uid-1", "alice@example.com", "viewer")

	if err := UpdateUserDisplayName(context.Background(), "uid-1", "Alice Smith"); err != nil {
		t.Fatalf("UpdateUserDisplayName failed: %v", err)
	}

	user, err := GetUserByUID(context.Background(), "uid-1")
	if err != nil {
		t.Fatalf("GetUserByUID failed: %v", err)
	}
	if user.DisplayName != "Alice Smith" {
		t.Errorf("expected display name 'Alice Smith', got '%s'", user.DisplayName)
	}
}

func TestDeleteRecipe_CascadesMakeLogs(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Cascade Log Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}
	if err := CreateMakeLog(ctx, &MakeLog{RecipeID: r.ID, MadeAt: "2024-01-15"}); err != nil {
		t.Fatalf("CreateMakeLog failed: %v", err)
	}

	if err := DeleteRecipe(ctx, r.ID); err != nil {
		t.Fatalf("DeleteRecipe failed: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM make_logs WHERE recipe_id = ?", r.ID).Scan(&count); err != nil {
		t.Fatalf("query make_logs failed: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 make_logs after recipe delete, got %d", count)
	}
}
