# Database Architecture

This document describes the SQLite + Cloud Storage database architecture for the RecipeBook application.

## Architecture Overview

```
┌─────────────────┐
│  React Frontend │
│  (Firebase Auth)│
└────────┬────────┘
         │ HTTP + JWT Token
         ▼
┌─────────────────────────┐
│   Cloud Run (Go API)    │
│  ┌──────────────────┐   │
│  │ SQLite (in /tmp) │   │
│  └──────────────────┘   │
└────────┬────────────────┘
         │ Download/Upload
         ▼
┌─────────────────────────┐
│  Cloud Storage Bucket   │
│    recipes.db (1-10MB)  │
│  (Versioning enabled)   │
└─────────────────────────┘
```

## Why SQLite + Cloud Storage?

### Cost Comparison
- **Cloud SQL**: $7-10/month for db-f1-micro
- **SQLite + Cloud Storage**: ~$0.001/month (essentially free)

### Benefits
1. **Full SQL support**: Complex queries, JOINs, full-text search
2. **Near-zero cost**: Storage is ~$0.02/GB/month, database is <10MB
3. **Simple deployment**: No database server to manage
4. **Full-text search**: Built-in FTS5 for ingredient/recipe search
5. **Markdown storage**: Simple TEXT fields for ingredients/method

### Trade-offs
- Not ideal for high concurrent writes (fine for 2 users)
- Eventual consistency (writes sync async to Cloud Storage)
- Container instances may have slightly stale data

## Database Schema

```sql
-- Main recipes table with UUID primary keys
CREATE TABLE recipes (
    id TEXT PRIMARY KEY,        -- UUID string (e.g., "550e8400-e29b-41d4-a716-446655440000")
    title TEXT NOT NULL,
    description TEXT,           -- Brief description of the recipe
    recipe_type TEXT,           -- 'food', 'cocktail', etc.
    cuisine TEXT,               -- 'italian', 'mexican', 'chinese', etc.
    ingredients TEXT,           -- markdown format
    method TEXT,                -- markdown format
    notes TEXT,                 -- markdown format
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tags table for recipe categorization
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL   -- lowercase tag names (e.g., 'pasta', 'quick', 'vegetarian')
);

-- Many-to-many relationship between recipes and tags
CREATE TABLE recipe_tags (
    recipe_id TEXT NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (recipe_id, tag_id),
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Recipe images stored in Google Cloud Storage
CREATE TABLE recipe_images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    recipe_id TEXT NOT NULL,
    image_url TEXT NOT NULL,    -- Full GCS URL (https://storage.googleapis.com/bucket/path)
    display_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
);
CREATE INDEX idx_recipe_images_recipe_id ON recipe_images(recipe_id);

-- Full-text search table (includes new description and cuisine fields)
CREATE VIRTUAL TABLE recipes_fts USING fts5(
    title, description, cuisine, ingredients, method, notes,
    content=recipes,
    content_rowid=id
);
```

**Access Control:** All recipes are public for reading. Role-based access control is implemented via Firestore:
- **Viewers**: Can read recipes (public access)
- **Editors**: Can create and update recipes
- **Admins**: Can delete recipes and manage user roles

**Image Storage:** Recipe images are stored in Google Cloud Storage, not in the database. The `recipe_images` table stores only the URLs and metadata.

## API Endpoints

### List all recipes
```
GET /recipes
(Public - no authentication required)

Response:
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Chicken Pasta",
    "description": "Creamy pasta with pan-fried chicken",
    "type": "food",
    "cuisine": "italian",
    "tags": ["pasta", "quick", "italian"],
    "ingredients": "## Ingredients\n\n- 2 chicken breasts\n- 500g pasta\n...",
    "method": "## Method\n\n1. Cook pasta...",
    "notes": "Optional notes...",
    "images": [
      {
        "id": 1,
        "recipeId": "550e8400-e29b-41d4-a716-446655440000",
        "imageUrl": "https://storage.googleapis.com/bucket/recipe-images/550e8400.../abc123.jpg",
        "displayOrder": 0
      }
    ],
    "createdAt": "2025-01-01T12:00:00Z",
    "updatedAt": "2025-01-01T12:00:00Z"
  }
]
```

### Get single recipe
```
GET /recipes/{id}
(Public - no authentication required)
```

### Create recipe
```
POST /recipes
Authorization: Bearer <firebase-token>
Content-Type: application/json

{
  "title": "Chicken Pasta",
  "description": "Creamy pasta with pan-fried chicken",
  "type": "food",
  "cuisine": "italian",
  "tags": ["pasta", "quick", "italian"],
  "ingredients": "## Ingredients\n\n- 2 chicken breasts\n- 500g pasta",
  "method": "## Method\n\n1. Cook pasta according to package",
  "notes": "Great for meal prep"
}
```

### Update recipe
```
PUT /recipes/{id}
Authorization: Bearer <firebase-token>
Content-Type: application/json

{
  "title": "Updated title",
  "description": "Updated description",
  "type": "food",
  "cuisine": "italian",
  "tags": ["pasta", "updated"],
  "ingredients": "...",
  "method": "...",
  "notes": "..."
}
```

### Delete recipe
```
DELETE /recipes/{id}
Authorization: Bearer <firebase-token>
(Requires admin role)
```

### Search recipes (full-text)
```
GET /recipes/search?q=chicken+pasta
(Public - no authentication required)

Searches across title, description, cuisine, ingredients, method, and notes fields.
```

### Upload recipe image
```
POST /recipes/images
Authorization: Bearer <firebase-token>
Content-Type: multipart/form-data

Form fields:
- recipeId: UUID of the recipe
- image: Image file (max 10MB, .jpg/.jpeg/.png/.gif/.webp)
- displayOrder: (optional) Integer for ordering multiple images

Response:
{
  "imageUrl": "https://storage.googleapis.com/bucket/recipe-images/550e8400.../abc123.jpg"
}
```

## Image Storage Architecture

Recipe images are stored separately from the database in Google Cloud Storage:

### Storage Structure
```
gs://{bucket-name}/
├── recipes.db                          # SQLite database file
└── recipe-images/                      # Image directory
    └── {recipe-uuid}/                  # One folder per recipe
        ├── {image-uuid-1}.jpg
        ├── {image-uuid-2}.png
        └── ...
```

### How Image Upload Works
1. User uploads image via multipart form to `/recipes/images` endpoint
2. Backend validates file size (max 10MB) and type (.jpg, .jpeg, .png, .gif, .webp)
3. Backend generates unique UUID filename: `recipe-images/{recipeId}/{uuid}.{ext}`
4. Image uploaded directly to GCS with public read access
5. Image URL and metadata saved to `recipe_images` table
6. Public URL returned: `https://storage.googleapis.com/{bucket}/{path}`

### Why This Approach?
- **Database stays small**: Images not stored in SQLite
- **Fast serving**: Images served directly from GCS (no backend involvement)
- **Simple**: Backend handles upload once, then forgets about it
- **Scalable**: GCS handles all image serving load
- **Cheap**: GCS storage is ~$0.02/GB/month

### Permissions
The same service account that accesses the database also handles image uploads. The `storage.objectAdmin` role on the bucket allows both database sync and image management.

## How It Works

### Startup
1. Cloud Run container starts
2. Downloads `recipes.db` from Cloud Storage to `/tmp/recipes.db`
3. If file doesn't exist, creates new database with schema
4. Opens SQLite connection

### Read Operations
- All reads use local SQLite file in `/tmp`
- Fast and efficient
- No network calls needed

### Write Operations
1. Write to local SQLite file immediately
2. Return response to client
3. Asynchronously upload entire database to Cloud Storage
4. Cloud Storage versioning keeps last 5 versions for safety

### Concurrent Writes
- For 2 users, concurrent writes are extremely rare
- Each Cloud Run instance has its own SQLite file
- Periodic sync keeps instances relatively in sync
- Last write wins (acceptable for personal use)

## Setup Instructions

### 1. Deploy Infrastructure

```bash
cd infra
terraform init
terraform plan
terraform apply
```

This creates:
- Cloud Storage bucket: `{project-id}-recipebook-db`
- Service account for Cloud Run with storage access
- Bucket versioning (keeps last 5 versions)

### 2. Update Backend Environment

When deploying to Cloud Run, set this environment variable:

```bash
DB_BUCKET_NAME={project-id}-recipebook-db
```

You can get the bucket name from Terraform output:
```bash
terraform output database_bucket_name
```

### 3. Local Development

For local development:

```bash
cd backend

# Set environment variables
export DB_BUCKET_NAME="{project-id}-recipebook-db"
export GOOGLE_APPLICATION_CREDENTIALS="./service-account.json"

# Run server
go run .
```

The database will be created locally in `/tmp/recipes.db` and synced to Cloud Storage.

### 4. Deploy to Cloud Run

```bash
cd backend

# Build and deploy
gcloud run deploy recipebook-backend \
  --source . \
  --region us-central1 \
  --service-account recipebook-backend@{project-id}.iam.gserviceaccount.com \
  --set-env-vars DB_BUCKET_NAME={project-id}-recipebook-db \
  --allow-unauthenticated
```

## Cost Breakdown

### Monthly Costs (Personal Use)

```
Cloud Storage:
- Database file: ~5 MB
- Cost: 5 MB × $0.02/GB = $0.0001/month

Operations:
- Class A (uploads): ~100/month × $0.05/10k = $0.0005
- Class B (downloads): ~100/month × $0.004/10k = $0.00004

Cloud Run:
- 10k requests/month: FREE (2M free tier)
- Compute: FREE (360k GB-seconds free tier)

Total: ~$0.001/month (essentially free)
```

### Comparison

| Solution | Monthly Cost | Search | Complexity |
|----------|-------------|--------|------------|
| Cloud SQL | $7-10 | Full SQL | Medium |
| SQLite + GCS | ~$0.001 | Full SQL + FTS | Medium |

## Backup & Recovery

### Automatic Backups
Cloud Storage versioning keeps the last 5 versions of the database automatically.

### Restore Previous Version
```bash
# List versions
gsutil ls -a gs://{bucket-name}/recipes.db

# Download specific version
gsutil cp gs://{bucket-name}/recipes.db#1234567890 ./recipes-backup.db
```

### Manual Backup
```bash
# Download current database
gsutil cp gs://{bucket-name}/recipes.db ./backup-$(date +%Y%m%d).db
```

## Monitoring

### Check Database Size
```bash
gsutil du -h gs://{bucket-name}/recipes.db
```

### View Upload Activity
```bash
gsutil logging get gs://{bucket-name}
```

### Cloud Run Logs
```bash
gcloud run services logs read recipebook-backend --region us-central1
```

## Limitations & Considerations

### When This Works Well ✅
- Personal or small team use (2-10 users)
- Read-heavy workloads
- Infrequent writes
- Need for full-text search
- Budget constraints

### When to Consider Alternatives ❌
- High concurrent writes (>10 users writing simultaneously)
- Need for strong consistency guarantees
- Real-time collaboration features
- Database size >100 MB
- Regulatory requirements for managed databases

## Future Enhancements

### Potential Improvements
1. **Read replicas**: Multiple Cloud Run instances download DB on startup
2. **Write queue**: Use Pub/Sub to coordinate writes across instances
3. **Caching**: Add Redis for frequently accessed recipes
4. **CDN**: Serve static recipe pages from CDN

### Migration to Cloud SQL
If you outgrow SQLite, migration to Cloud SQL is straightforward:
1. Export SQLite to SQL dump: `sqlite3 recipes.db .dump > dump.sql`
2. Import to PostgreSQL/MySQL
3. Update connection string in code
4. No schema changes needed

## Troubleshooting

### Database not syncing
- Check service account permissions on bucket
- Check Cloud Run logs for upload errors
- Verify `DB_BUCKET_NAME` environment variable

### Old data showing
- Cloud Run instances cache the DB file
- Force refresh by redeploying service
- Or wait for new container instances to spin up

### Database corruption
- Restore from Cloud Storage version
- SQLite WAL mode prevents most corruption
- Versioning keeps 5 previous versions

## Questions?

This architecture is optimized for:
- **Cost**: Nearly free
- **Simplicity**: No database server
- **Features**: Full SQL + full-text search
- **Use case**: Personal/small team recipe management

For questions or issues, see the main README.md
