package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Recipe struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"` // "food", "cocktail", etc.
	Ingredients []string `json:"ingredients"`
	Steps       []string `json:"steps"`
}

func main() {
	ctx := context.Background()

	// Initialize Firebase Admin SDK
	if err := InitFirebase(ctx); err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", healthHandler)
	// Protected route - requires authentication
	http.HandleFunc("/recipes", AuthMiddleware(recipesHandler))

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

	// Get authenticated user ID from context
	userID := r.Context().Value("userID").(string)
	log.Printf("Request from user: %s", userID)

	switch r.Method {
	case http.MethodGet:
		// TODO: Implement get recipes from Firestore
		json.NewEncoder(w).Encode([]Recipe{})
	case http.MethodPost:
		// TODO: Implement create recipe in Firestore
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Recipe created"})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
