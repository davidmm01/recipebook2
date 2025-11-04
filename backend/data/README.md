# Recipe Import Tool

This directory contains YAML recipe files and a tool to import them into the database.

## Directory Structure

- `cocktails/` - Contains YAML files for drink recipes
- `recipes/` - Contains YAML files for food recipes
- `instructions` - Instructions for using the import tool

## Quick Start

**IMPORTANT:** The backend must be running before importing recipes. The backend initializes the database with the correct schema on startup.

From the `backend` directory:

### Step 1: Start the backend
```bash
make run-local
```

Wait for the backend to start and create the database (you'll see "Database initialized successfully" in the logs).

### Step 2: Import recipes (in a new terminal)
```bash
# Generate fresh SQL import file
make import-recipes

# Import into the running database
sqlite3 /tmp/recipes.db < import.sql
```

Or in one command:
```bash
make import-and-load
```

### Alternative: Manual import
```bash
# Generate SQL import file
go run cmd/import-recipes/main.go .

# Import into database
sqlite3 /tmp/recipes.db < import.sql
```

### Why this order matters

The backend's `InitDatabase` function downloads the database from Cloud Storage, or creates a new empty one if none exists. If you import recipes before starting the backend, they will be overwritten when the backend initializes. Therefore, always start the backend first, then import.

## What the Import Tool Does

The import script (`cmd/import-recipes/main.go`) converts YAML recipe files into SQL INSERT statements for the database. It handles the following mappings:

### Field Mappings

| Database Field | Source | Notes |
|---------------|---------|-------|
| `id` | Generated | UUID v4 |
| `title` | `name` from YAML | |
| `description` | Empty | New field, not in YAML |
| `recipe_type` | Directory | "food" for recipes/, "drink" for cocktails/ |
| `cuisine` | `cuisine` from YAML | |
| `tags` | `descriptors` from YAML | Stored in tags table |
| `ingredients` | `ingredients` from YAML | Converted to markdown with subheadings |
| `method` | `instructions` from YAML | Converted to markdown with subheadings |
| `notes` | `notes` + `next` from YAML | Combined with subheadings |
| `sources` | `source` fields from YAML | Excludes submitter field |
| `created_at` | `date_added` from YAML | |
| `updated_at` | Current time | |
| `created_by_name` | `source.submitter` from YAML | |

### Special Features

#### Comment to Subheading Conversion

YAML comments in ingredient or instruction lists are converted to markdown subheadings. For example:

**YAML:**
```yaml
ingredients:
  - 500g pork mince
  - 900g tofu
  # sauces
  - 113g mapo tofu sauce
  - 1 tablespoon oyster sauce
```

**Converts to:**
```markdown
- 500g pork mince
- 900g tofu

### sauces
- 113g mapo tofu sauce
- 1 tablespoon oyster sauce
```

#### Notes Combination

The `notes` and `next` fields from YAML are combined into a single notes field with subheadings:

**YAML:**
```yaml
notes:
  - Serve with rice
  - Great with pickled vegetables
next:
  - Try with lamb instead of beef
  - Add more vegetables
```

**Converts to:**
```markdown
### Notes
- Serve with rice
- Great with pickled vegetables

### Next
- Try with lamb instead of beef
- Add more vegetables
```

## Output

The script generates an `import.sql` file containing:
- SQL INSERT statements for all recipes
- INSERT statements for tags
- Relationship records linking recipes to tags
- All statements wrapped in a transaction

## Adding New Recipes

1. Create a new YAML file in `cocktails/` or `recipes/` directory
2. Follow the template structure (see `recipes/template.yaml`)
3. Run `make import-recipes` to regenerate the SQL file
4. Import into your database with `sqlite3 /tmp/recipes.db < import.sql`

## Notes

- The script skips `template.yaml` files automatically
- Empty or malformed YAML files will generate warnings but won't stop the import
- All tags are automatically converted to lowercase
- SQL strings are properly escaped to prevent injection issues
- The script is idempotent - you can run it multiple times safely
