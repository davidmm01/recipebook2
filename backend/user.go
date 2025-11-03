package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// userProfileHandler handles GET and PUT requests for user profiles
func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Auth required
	userID, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get user's own profile
		// Note: authenticateRequest above already created the user if they didn't exist
		user, err := GetUserByUID(r.Context(), userID)
		if err != nil {
			log.Printf("Error getting user profile for %s: %v", userID, err)
			http.Error(w, fmt.Sprintf("Failed to get user profile: %v", err), http.StatusInternalServerError)
			return
		}

		if user == nil {
			log.Printf("User %s not found in database after authentication", userID)
			http.Error(w, "User not found - this shouldn't happen after authentication", http.StatusNotFound)
			return
		}

		log.Printf("Successfully retrieved profile for user %s (%s)", userID, user.Email)
		json.NewEncoder(w).Encode(user)

	case http.MethodPut:
		// Update user's own profile
		var updates map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Only allow updating displayName
		displayName, ok := updates["displayName"]
		if !ok {
			http.Error(w, "No displayName field provided", http.StatusBadRequest)
			return
		}

		displayNameStr, ok := displayName.(string)
		if !ok || displayNameStr == "" {
			http.Error(w, "Invalid displayName value", http.StatusBadRequest)
			return
		}

		if err := UpdateUserDisplayName(r.Context(), userID, displayNameStr); err != nil {
			log.Printf("Error updating user display name: %v", err)
			http.Error(w, "Failed to update user profile", http.StatusInternalServerError)
			return
		}

		// Return updated profile
		user, err := GetUserByUID(r.Context(), userID)
		if err != nil {
			log.Printf("Error getting updated user profile: %v", err)
			http.Error(w, "Failed to get updated profile", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(user)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
