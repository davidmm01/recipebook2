package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	_ "github.com/mattn/go-sqlite3"
)

const (
	localDBPath = "/tmp/recipes.db"
)

type User struct {
	FirebaseUID string
	Email       string
	DisplayName string
	Role        string
}

func main() {
	// Define commands
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	setRoleCmd := flag.NewFlagSet("set-role", flag.ExitOnError)

	// set-role flags
	email := setRoleCmd.String("email", "", "User email address")
	role := setRoleCmd.String("role", "", "New role (viewer, editor, or admin)")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		if err := listUsers(); err != nil {
			log.Fatalf("Error: %v", err)
		}

	case "set-role":
		setRoleCmd.Parse(os.Args[2:])
		if *email == "" || *role == "" {
			fmt.Println("Error: Both --email and --role are required")
			setRoleCmd.PrintDefaults()
			os.Exit(1)
		}

		if *role != "viewer" && *role != "editor" && *role != "admin" {
			fmt.Println("Error: Role must be one of: viewer, editor, admin")
			os.Exit(1)
		}

		if err := setUserRole(*email, *role); err != nil {
			log.Fatalf("Error: %v", err)
		}

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("User Management CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  manage-users list")
	fmt.Println("  manage-users set-role --email <email> --role <viewer|editor|admin>")
	fmt.Println("\nExamples:")
	fmt.Println("  manage-users list")
	fmt.Println("  manage-users set-role --email user@example.com --role editor")
	fmt.Println("  manage-users set-role --email admin@example.com --role admin")
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", localDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func listUsers() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT firebase_uid, email, display_name, role
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		return fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		var displayName sql.NullString
		if err := rows.Scan(&u.FirebaseUID, &u.Email, &displayName, &u.Role); err != nil {
			return fmt.Errorf("failed to scan user: %w", err)
		}
		u.DisplayName = displayName.String
		users = append(users, u)
	}

	if len(users) == 0 {
		fmt.Println("No users found in database")
		return nil
	}

	// Print users in a table
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "EMAIL\tDISPLAY NAME\tROLE\tFIREBASE UID")
	fmt.Fprintln(w, "-----\t------------\t----\t------------")

	for _, u := range users {
		displayName := u.DisplayName
		if displayName == "" {
			displayName = "(not set)"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", u.Email, displayName, u.Role, u.FirebaseUID[:12]+"...")
	}
	w.Flush()

	fmt.Printf("\nTotal users: %d\n", len(users))
	return nil
}

func setUserRole(email, newRole string) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	ctx := context.Background()

	// Check if user exists
	var currentRole string
	var displayName sql.NullString
	err = db.QueryRowContext(ctx, "SELECT role, display_name FROM users WHERE email = ?", email).
		Scan(&currentRole, &displayName)

	if err == sql.ErrNoRows {
		return fmt.Errorf("user not found with email: %s", email)
	}
	if err != nil {
		return fmt.Errorf("failed to query user: %w", err)
	}

	// Update role
	_, err = db.ExecContext(ctx, "UPDATE users SET role = ? WHERE email = ?", newRole, email)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	name := displayName.String
	if name == "" {
		name = email
	}

	fmt.Printf("✓ Successfully updated %s's role from '%s' to '%s'\n", name, currentRole, newRole)
	return nil
}
