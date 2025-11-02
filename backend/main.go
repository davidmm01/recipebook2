package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	// Filter metadata endpoints
	http.HandleFunc("/tags", corsMiddleware(tagsHandler))
	http.HandleFunc("/cuisines", corsMiddleware(cuisinesHandler))
	// Image upload endpoint
	http.HandleFunc("/recipes/images", corsMiddleware(imageUploadHandler))
	// Icon endpoints
	http.HandleFunc("/icons", corsMiddleware(iconsHandler))

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
		// Check for filter parameters
		searchQuery := r.URL.Query().Get("search")
		cuisine := r.URL.Query().Get("cuisine")
		tagsParam := r.URL.Query()["tags"] // Get all tags parameters (can be multiple)

		var recipes []Recipe
		var err error

		// If any filters are provided, use FilterRecipes
		if searchQuery != "" || cuisine != "" || len(tagsParam) > 0 {
			recipes, err = FilterRecipes(r.Context(), searchQuery, tagsParam, cuisine)
		} else {
			// No filters, get all recipes
			recipes, err = GetRecipes(r.Context())
		}

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

	// Extract recipe ID from URL path (now a UUID string)
	recipeID := strings.TrimPrefix(r.URL.Path, "/recipes/")
	if recipeID == "" || recipeID == "search" {
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

func tagsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Public read - no auth required
	tags, err := GetAllTags(r.Context())
	if err != nil {
		log.Printf("Error getting tags: %v", err)
		http.Error(w, "Failed to get tags", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tags)
}

func cuisinesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Public read - no auth required
	cuisines, err := GetAllCuisines(r.Context())
	if err != nil {
		log.Printf("Error getting cuisines: %v", err)
		http.Error(w, "Failed to get cuisines", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cuisines)
}

func imageUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Auth required
	userID, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("Uploading image - authenticated user: %s", userID)

	// Parse multipart form (limit to 10MB)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get recipe ID
	recipeID := r.FormValue("recipeId")
	if recipeID == "" {
		http.Error(w, "Missing recipeId", http.StatusBadRequest)
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get image from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Upload to GCS
	imageURL, err := UploadImageToGCS(r.Context(), file, fileHeader, recipeID)
	if err != nil {
		log.Printf("Error uploading image: %v", err)
		http.Error(w, fmt.Sprintf("Failed to upload image: %v", err), http.StatusInternalServerError)
		return
	}

	// Get display order (optional, defaults to 0)
	displayOrder := 0
	if orderStr := r.FormValue("displayOrder"); orderStr != "" {
		if order, err := fmt.Sscanf(orderStr, "%d", &displayOrder); err == nil && order == 1 {
			// Successfully parsed
		}
	}

	// Save to database
	dbMutex.Lock()
	err = addRecipeImage(r.Context(), recipeID, imageURL, displayOrder)
	dbMutex.Unlock()

	if err != nil {
		log.Printf("Error saving image to database: %v", err)
		// Try to delete from GCS
		DeleteImageFromGCS(r.Context(), imageURL)
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	// Return the image URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"imageUrl": imageURL,
	})
}

func iconsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Public read - no auth required
		icons, err := GetAllIcons(r.Context())
		if err != nil {
			log.Printf("Error getting icons: %v", err)
			http.Error(w, "Failed to get icons", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(icons)

	case http.MethodPost:
		// Auth required for uploads
		userID, err := authenticateRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("Uploading icon - authenticated user: %s", userID)

		// Parse multipart form (limit to 2MB for icons)
		err = r.ParseMultipartForm(2 << 20)
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		// Get file from form
		file, fileHeader, err := r.FormFile("icon")
		if err != nil {
			http.Error(w, "Failed to get icon from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Upload to GCS
		filename, iconURL, err := UploadIconToGCS(r.Context(), file, fileHeader)
		if err != nil {
			log.Printf("Error uploading icon: %v", err)
			http.Error(w, fmt.Sprintf("Failed to upload icon: %v", err), http.StatusInternalServerError)
			return
		}

		// Save to database
		dbMutex.Lock()
		iconID, err := createIcon(r.Context(), filename, iconURL)
		dbMutex.Unlock()

		if err != nil {
			log.Printf("Error saving icon to database: %v", err)
			// Try to delete from GCS
			DeleteImageFromGCS(r.Context(), iconURL)
			http.Error(w, "Failed to save icon", http.StatusInternalServerError)
			return
		}

		// Return the created icon
		icon := Icon{
			ID:       iconID,
			Filename: filename,
			IconURL:  iconURL,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(icon)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
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
