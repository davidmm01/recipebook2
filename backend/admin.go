package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// requireAdmin authenticates the request and verifies the caller has the admin role.
// Returns the caller's DBUser on success, or writes an HTTP error and returns nil.
func requireAdmin(w http.ResponseWriter, r *http.Request) *DBUser {
	userID, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}
	user, err := GetUserByUID(r.Context(), userID)
	if err != nil || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}
	if user.Role != "admin" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return nil
	}
	return user
}

// adminUsersListHandler handles GET /admin/users — returns all users (admin only).
func adminUsersListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if requireAdmin(w, r) == nil {
		return
	}

	users, err := GetAllUsers(r.Context())
	if err != nil {
		log.Printf("Error listing users: %v", err)
		http.Error(w, "Failed to list users", http.StatusInternalServerError)
		return
	}

	if users == nil {
		users = []DBUser{}
	}
	json.NewEncoder(w).Encode(users)
}

// adminUpdateRoleHandler handles PUT /admin/users/{uid}/role — updates a user's role (admin only).
func adminUpdateRoleHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	caller := requireAdmin(w, r)
	if caller == nil {
		return
	}

	// Extract UID from path: /admin/users/{uid}/role
	uid := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/admin/users/"), "/role")
	if uid == "" {
		http.Error(w, "Missing user UID", http.StatusBadRequest)
		return
	}

	if uid == caller.FirebaseUID {
		http.Error(w, "Cannot change your own role", http.StatusForbidden)
		return
	}

	var body struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	validRoles := map[string]bool{"viewer": true, "editor": true, "admin": true}
	if !validRoles[body.Role] {
		http.Error(w, "Invalid role: must be viewer, editor, or admin", http.StatusBadRequest)
		return
	}

	if err := UpdateUserRole(r.Context(), uid, body.Role); err != nil {
		log.Printf("Error updating role for user %s: %v", uid, err)
		http.Error(w, "Failed to update user role", http.StatusInternalServerError)
		return
	}

	user, err := GetUserByUID(r.Context(), uid)
	if err != nil || user == nil {
		http.Error(w, "Failed to retrieve updated user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}
