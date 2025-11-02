package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	ctx := context.Background()

	// Initialize Firebase Admin SDK for authentication
	if err := InitFirebase(ctx); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize SQLite database
	bucketName := os.Getenv("DB_BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("DB_BUCKET_NAME environment variable is required")
	}

	if err := InitDatabase(ctx, bucketName); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", corsMiddleware(healthHandler))
	// Public read, auth required for writes
	http.HandleFunc("/recipes", corsMiddleware(recipesHandler))
	http.HandleFunc("/recipes/", corsMiddleware(recipeByIDHandler))
	http.HandleFunc("/recipes/search", corsMiddleware(searchHandler))

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func recipesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Public read - no auth required
		recipes, err := GetRecipes(r.Context())
		if err != nil {
			log.Printf("Error getting recipes: %v", err)
			http.Error(w, "Failed to get recipes", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(recipes)

	case http.MethodPost:
		// Auth required for writes
		userID, err := authenticateRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("Creating recipe - authenticated user: %s", userID)

		var recipe Recipe
		if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := CreateRecipe(r.Context(), &recipe); err != nil {
			log.Printf("Error creating recipe: %v", err)
			http.Error(w, "Failed to create recipe", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(recipe)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func recipeByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract recipe ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/recipes/")
	recipeID, err := strconv.ParseInt(path, 10, 64)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Public read - no auth required
		recipe, err := GetRecipeByID(r.Context(), recipeID)
		if err != nil {
			log.Printf("Error getting recipe: %v", err)
			http.Error(w, "Failed to get recipe", http.StatusInternalServerError)
			return
		}
		if recipe == nil {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(recipe)

	case http.MethodPut:
		// Auth required for writes
		userID, err := authenticateRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("Updating recipe - authenticated user: %s", userID)

		var recipe Recipe
		if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		recipe.ID = recipeID

		if err := UpdateRecipe(r.Context(), &recipe); err != nil {
			log.Printf("Error updating recipe: %v", err)
			http.Error(w, "Failed to update recipe", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(recipe)

	case http.MethodDelete:
		// Auth required for writes
		userID, err := authenticateRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("Deleting recipe - authenticated user: %s", userID)

		if err := DeleteRecipe(r.Context(), recipeID); err != nil {
			log.Printf("Error deleting recipe: %v", err)
			http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Public read - no auth required
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing search query parameter 'q'", http.StatusBadRequest)
		return
	}

	recipes, err := SearchRecipes(r.Context(), query)
	if err != nil {
		log.Printf("Error searching recipes: %v", err)
		http.Error(w, "Failed to search recipes", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(recipes)
}

// corsMiddleware adds CORS headers to allow cross-origin requests from the frontend
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from frontend (localhost:3000 for local dev)
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:3000" || origin == "" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

// authenticateRequest validates Firebase ID token and returns user ID
func authenticateRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	idToken := parts[1]
	token, err := firebaseAuth.VerifyIDToken(r.Context(), idToken)
	if err != nil {
		return "", fmt.Errorf("invalid or expired token: %w", err)
	}

	return token.UID, nil
}
