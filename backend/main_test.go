package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Health handler ---

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", body["status"])
	}
}

// --- CORS middleware ---

func TestCORSMiddleware_SetsHeaders(t *testing.T) {
	handler := corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("expected Access-Control-Allow-Origin header to be set")
	}
}

func TestCORSMiddleware_PreflightOptions(t *testing.T) {
	handler := corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		// Should not be called for OPTIONS preflight
		t.Error("inner handler should not be called for OPTIONS request")
	})

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204 for OPTIONS preflight, got %d", w.Code)
	}
}

// --- Recipes handler ---

func TestRecipesHandler_GetEmpty(t *testing.T) {
	clearTables(t)

	req := httptest.NewRequest(http.MethodGet, "/recipes", nil)
	w := httptest.NewRecorder()

	recipesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRecipesHandler_GetWithData(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Handler Test Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/recipes", nil)
	w := httptest.NewRecorder()

	recipesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var recipes []Recipe
	if err := json.NewDecoder(w.Body).Decode(&recipes); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(recipes) != 1 {
		t.Errorf("expected 1 recipe, got %d", len(recipes))
	}
	if len(recipes) > 0 && recipes[0].Title != "Handler Test Recipe" {
		t.Errorf("expected title 'Handler Test Recipe', got '%s'", recipes[0].Title)
	}
}

func TestRecipesHandler_PostUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/recipes", nil)
	w := httptest.NewRecorder()

	recipesHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// --- Recipe by ID handler ---

func TestRecipeByIDHandler_GetFound(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Find By ID")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/recipes/"+r.ID, nil)
	w := httptest.NewRecorder()

	recipeByIDHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var got Recipe
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Title != "Find By ID" {
		t.Errorf("expected title 'Find By ID', got '%s'", got.Title)
	}
}

func TestRecipeByIDHandler_GetNotFound(t *testing.T) {
	clearTables(t)

	req := httptest.NewRequest(http.MethodGet, "/recipes/nonexistent-id", nil)
	w := httptest.NewRecorder()

	recipeByIDHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestRecipeByIDHandler_PutUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/recipes/some-id", nil)
	w := httptest.NewRecorder()

	recipeByIDHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestRecipeByIDHandler_DeleteUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/recipes/some-id", nil)
	w := httptest.NewRecorder()

	recipeByIDHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

// --- Search handler ---

func TestSearchHandler_MissingQuery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/recipes/search", nil)
	w := httptest.NewRecorder()

	searchHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing query, got %d", w.Code)
	}
}

// --- Tags handler ---

func TestTagsHandler_Get(t *testing.T) {
	clearTables(t)

	req := httptest.NewRequest(http.MethodGet, "/tags", nil)
	w := httptest.NewRecorder()

	tagsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// --- Cuisines handler ---

func TestCuisinesHandler_Get(t *testing.T) {
	clearTables(t)

	req := httptest.NewRequest(http.MethodGet, "/cuisines", nil)
	w := httptest.NewRecorder()

	cuisinesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// --- Make logs handler ---

func TestMakeLogsHandler_GetEmpty(t *testing.T) {
	clearTables(t)
	ctx := context.Background()

	r := makeTestRecipe("Make Log Handler Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/make-logs/"+r.ID, nil)
	w := httptest.NewRecorder()

	makeLogsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMakeLogsHandler_PostUnauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/make-logs/some-recipe-id", nil)
	w := httptest.NewRecorder()

	makeLogsHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
