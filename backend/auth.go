package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var firebaseAuth *auth.Client

// InitFirebase initializes Firebase Admin SDK
func InitFirebase(ctx context.Context, configPath string) error {
	// If running in GCP (Cloud Run), Firebase will use Application Default Credentials
	// For local development, set firebase_service_account_path in config or GOOGLE_APPLICATION_CREDENTIALS env var
	var app *firebase.App
	var err error

	if path := serviceAccountPath(configPath); path != "" {
		app, err = firebase.NewApp(ctx, nil, option.WithCredentialsFile(path))
	} else {
		// Use default credentials (works in Cloud Run)
		app, err = firebase.NewApp(ctx, nil)
	}

	if err != nil {
		return err
	}

	firebaseAuth, err = app.Auth(ctx)
	if err != nil {
		return err
	}

	log.Println("Firebase Admin SDK initialized")
	return nil
}

func serviceAccountPath(configPath string) string {
	// Check config file path first
	if configPath != "" {
		return configPath
	}
	// Fall back to environment variable
	if path := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"); path != "" {
		return path
	}
	return ""
}

// AuthMiddleware validates Firebase ID token from Authorization header
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		idToken := parts[1]
		token, err := firebaseAuth.VerifyIDToken(r.Context(), idToken)
		if err != nil {
			log.Printf("Error verifying token: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context for use in handlers
		ctx := context.WithValue(r.Context(), "userID", token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
