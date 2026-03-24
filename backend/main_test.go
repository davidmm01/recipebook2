package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// withAuth stubs the authenticate function to return the given UID, bypassing Firebase.
// Call as: defer withAuth("uid")()
func withAuth(uid string) func() {
	original := authenticate
	authenticate = func(r *http.Request) (string, error) {
		return uid, nil
	}
	return func() { authenticate = original }
}

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

// --- Admin handlers ---

func TestAdminUsersListHandler_Unauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	w := httptest.NewRecorder()

	adminUsersListHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAdminUsersListHandler_MethodNotAllowed(t *testing.T) {
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/admin/users", nil)
		w := httptest.NewRecorder()

		adminUsersListHandler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected 405, got %d", method, w.Code)
		}
	}
}

func TestAdminUpdateRoleHandler_Unauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/admin/users/some-uid/role", nil)
	w := httptest.NewRecorder()

	adminUpdateRoleHandler(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAdminUpdateRoleHandler_MethodNotAllowed(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodDelete} {
		req := httptest.NewRequest(method, "/admin/users/some-uid/role", nil)
		w := httptest.NewRecorder()

		adminUpdateRoleHandler(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("method %s: expected 405, got %d", method, w.Code)
		}
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

// --- Tests enabled by injectable auth ---

func TestRecipesHandler_PostCreatesRecipe(t *testing.T) {
	clearTables(t)
	defer withAuth("user-uid-1")()
	makeTestUser("user-uid-1", "alice@example.com", "editor")

	body := `{"title":"New Recipe","description":"desc","recipeType":"food","cuisine":"italian","ingredients":"flour","method":"mix"}`
	req := httptest.NewRequest(http.MethodPost, "/recipes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	recipesHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var got Recipe
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Title != "New Recipe" {
		t.Errorf("expected title 'New Recipe', got '%s'", got.Title)
	}
	if got.CreatedByUserID == nil || *got.CreatedByUserID != "user-uid-1" {
		t.Error("expected CreatedByUserID to be set to the authenticated user's UID")
	}
}

func TestRecipeByIDHandler_DeleteAsOwner(t *testing.T) {
	clearTables(t)
	defer withAuth("owner-uid")()
	ctx := context.Background()
	makeTestUser("owner-uid", "owner@example.com", "editor")

	ownerUID := "owner-uid"
	r := makeTestRecipe("Delete Me")
	r.CreatedByUserID = &ownerUID
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/recipes/"+r.ID, nil)
	w := httptest.NewRecorder()

	recipeByIDHandler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d: %s", w.Code, w.Body.String())
	}

	got, err := GetRecipeByID(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	if got != nil {
		t.Error("expected recipe to be deleted")
	}
}

func TestRecipeByIDHandler_DeleteAsAdmin(t *testing.T) {
	clearTables(t)
	defer withAuth("admin-uid")()
	ctx := context.Background()
	makeTestUser("admin-uid", "admin@example.com", "admin")
	makeTestUser("other-uid", "other@example.com", "editor")

	otherUID := "other-uid"
	r := makeTestRecipe("Admin Can Delete This")
	r.CreatedByUserID = &otherUID
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/recipes/"+r.ID, nil)
	w := httptest.NewRecorder()

	recipeByIDHandler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRecipeByIDHandler_DeleteForbiddenForNonOwner(t *testing.T) {
	clearTables(t)
	defer withAuth("viewer-uid")()
	ctx := context.Background()
	makeTestUser("viewer-uid", "viewer@example.com", "viewer")
	makeTestUser("other-uid", "other@example.com", "editor")

	otherUID := "other-uid"
	r := makeTestRecipe("Not Your Recipe")
	r.CreatedByUserID = &otherUID
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/recipes/"+r.ID, nil)
	w := httptest.NewRecorder()

	recipeByIDHandler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestRecipeByIDHandler_PutUpdatesRecipe(t *testing.T) {
	clearTables(t)
	defer withAuth("user-uid-1")()
	ctx := context.Background()

	r := makeTestRecipe("Original")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	body := `{"title":"Updated","description":"desc","recipeType":"food","cuisine":"french","ingredients":"butter","method":"melt"}`
	req := httptest.NewRequest(http.MethodPut, "/recipes/"+r.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	recipeByIDHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	got, err := GetRecipeByID(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetRecipeByID failed: %v", err)
	}
	if got.Title != "Updated" {
		t.Errorf("expected title 'Updated', got '%s'", got.Title)
	}
}

func TestMakeLogsHandler_PostCreatesMakeLog(t *testing.T) {
	clearTables(t)
	defer withAuth("user-uid-1")()
	ctx := context.Background()
	makeTestUser("user-uid-1", "alice@example.com", "editor")

	r := makeTestRecipe("Logged Recipe")
	if err := CreateRecipe(ctx, r); err != nil {
		t.Fatalf("CreateRecipe failed: %v", err)
	}

	body := `{"madeAt":"2024-06-01","notes":"Great result"}`
	req := httptest.NewRequest(http.MethodPost, "/make-logs/"+r.ID, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	makeLogsHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	logs, err := GetMakeLogsByRecipe(ctx, r.ID)
	if err != nil {
		t.Fatalf("GetMakeLogsByRecipe failed: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("expected 1 make log, got %d", len(logs))
	}
}

func TestUserProfileHandler_GetAuthorized(t *testing.T) {
	clearTables(t)
	defer withAuth("user-uid-1")()
	makeTestUser("user-uid-1", "alice@example.com", "editor")

	req := httptest.NewRequest(http.MethodGet, "/user/profile", nil)
	w := httptest.NewRecorder()

	userProfileHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var got DBUser
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Email != "alice@example.com" {
		t.Errorf("expected email 'alice@example.com', got '%s'", got.Email)
	}
	if got.Role != "editor" {
		t.Errorf("expected role 'editor', got '%s'", got.Role)
	}
}
