# User Management CLI

A command-line tool for managing user roles in the RecipeBook application.

## Prerequisites

- Go 1.16 or higher
- Access to the SQLite database (default: `/tmp/recipes.db`)

## Installation

From the `backend` directory:

```bash
go build -o manage-users ./cmd/manage-users
```

This creates an executable called `manage-users` in the current directory.

## Usage

### List all users

```bash
./manage-users list
```

Example output:
```
EMAIL                    DISPLAY NAME    ROLE      FIREBASE UID
-----                    ------------    ----      ------------
alice@example.com        Alice Smith     admin     abc123def456...
bob@example.com          Bob Jones       editor    xyz789uvw012...
charlie@example.com      (not set)       viewer    mno345pqr678...

Total users: 3
```

### Change a user's role

```bash
./manage-users set-role --email <email> --role <viewer|editor|admin>
```

**Examples:**

Make a user an editor:
```bash
./manage-users set-role --email bob@example.com --role editor
```

Promote a user to admin:
```bash
./manage-users set-role --email alice@example.com --role admin
```

Demote a user to viewer:
```bash
./manage-users set-role --email charlie@example.com --role viewer
```

## Role Hierarchy

- **viewer**: Can only read recipes (default for new users)
- **editor**: Can create and update recipes
- **admin**: Can delete recipes and manage users

## Notes

- Users are automatically created with the "viewer" role on first login
- The database must exist at `/tmp/recipes.db` (or update `localDBPath` in the code)
- This tool directly modifies the SQLite database
- Changes take effect immediately - no server restart required
