# Cloud Run Database Behavior

This document explains how the SQLite + Cloud Storage pattern works inside Cloud Run and the eventual consistency considerations.

## How It Works

### Container Lifecycle

```
Cloud Run Container Starts
    ↓
1. Downloads recipes.db from Cloud Storage → /tmp/recipes.db
    ↓
2. Opens SQLite connection to /tmp/recipes.db
    ↓
3. Container is ready to handle requests
    ↓
4. Reads: Query local SQLite file (fast, no network)
    ↓
5. Writes: Update local SQLite + async upload to Cloud Storage
    ↓
Container Shutdown (after ~15 min idle timeout)
```

### The /tmp Directory

- Each Cloud Run container instance gets its own ephemeral `/tmp` directory
- Stored in memory (fast)
- **Disappears when container shuts down**
- **Survives between requests while container is alive**
- Typically ~512MB available

### Database Download Pattern

```go
func init() {
    // This happens ONCE per container startup
    downloadDBFromGCS()  // Gets latest recipes.db
    db, _ = sql.Open("sqlite3", "/tmp/recipes.db")
}
```

The database is **only downloaded on container startup**, not on every request.

---

## Multiple Container Instances

Cloud Run can spin up multiple containers when traffic increases:

```
Request 1 → Container A (/tmp/recipes.db - version at T=0)
Request 2 → Container A (reuses same container)
Request 3 → Container B (/tmp/recipes.db - version at T=10) <-- New container!
```

Each container:
- Downloads the database **once on startup**
- Keeps it in `/tmp` for its lifetime (~15 min idle timeout)
- Uploads to Cloud Storage after every write
- **Does NOT re-download after writes**

---

## Eventual Consistency Problem

### Example Scenario

```
Timeline:

T=0s:   Container A starts, downloads DB (10 recipes)
T=10s:  You create recipe #11 in Container A
T=11s:  Container A uploads DB to Cloud Storage ✅
T=15s:  Request hits Container A → sees 11 recipes ✅

T=30s:  Container B starts (traffic spike), downloads DB (11 recipes)
T=45s:  You create recipe #12 in Container B
T=46s:  Container B uploads DB to Cloud Storage ✅

T=50s:  Girlfriend's request hits Container A (still has old /tmp)
        → sees only 11 recipes ❌ (missing recipe #12)

T=65s:  Container A times out (idle for 15 min)
T=70s:  New request → Container C starts, downloads DB (12 recipes)
        → sees 12 recipes ✅
```

**The Problem:** Container A doesn't know Container B created a new recipe. Each container's `/tmp/recipes.db` is independent.

---

## Solutions

### ✅ **RECOMMENDED: Single Instance (maxScale = 1)**

Force Cloud Run to only ever run one container:

```hcl
# In infra/main.tf - when adding Cloud Run service
resource "google_cloud_run_service" "backend" {
  # ... other config ...

  template {
    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale" = "0"  # Scale to zero when idle
        "autoscaling.knative.dev/maxScale" = "1"  # Only 1 container max
      }
    }
  }
}
```

**Why This Works:**
- Only one container = only one `/tmp/recipes.db`
- No synchronization issues
- SQLite handles concurrent requests fine (multiple readers, single writer)
- Still scales to zero (free when idle)
- One container is plenty for 2 users

**Performance:**
- Single container can handle ~100+ requests/second
- Way more than enough for 2 users
- No added latency
- Still free (within Cloud Run free tier)

**This is the simplest and best solution for your use case.**

---

### Option 2: Periodic Re-download (More Complex)

Add background sync to refresh database periodically:

```go
func init() {
    downloadDBFromGCS()
    db, _ = sql.Open("sqlite3", localDBPath)

    // Re-download every 5 minutes
    go func() {
        ticker := time.NewTicker(5 * time.Minute)
        for range ticker.C {
            refreshDatabase()
        }
    }()
}

func refreshDatabase() {
    dbMutex.Lock()
    defer dbMutex.Unlock()

    // Download latest
    downloadDBFromGCS()

    // Reopen connection
    db.Close()
    db, _ = sql.Open("sqlite3", localDBPath)

    log.Println("Database refreshed from Cloud Storage")
}
```

**Pros:**
- Containers auto-refresh
- Max staleness: 5 minutes

**Cons:**
- More complex
- Extra Cloud Storage downloads (minimal cost)
- Still eventual consistency
- Not needed if you use single instance

---

### Option 3: Check Cloud Storage Before Each Read (Overkill)

Check if Cloud Storage has a newer version before every read:

```go
func GetRecipes(ctx context.Context) ([]Recipe, error) {
    // Check Cloud Storage modification time
    if cloudStorageIsNewer() {
        refreshDatabase()
    }

    // Query local SQLite
    return queryRecipes()
}
```

**Pros:**
- Always up-to-date

**Cons:**
- Extra network call per request (~50-100ms latency)
- More complex
- Unnecessary for 2 users
- Don't do this

---

## Real-World Behavior (Current Implementation)

### Typical Usage Pattern for 2 Users

```
Morning (8am):
- You open app → New container starts, downloads latest DB
- You browse recipes (same container handles all requests)
- You add a recipe → Uploads to Cloud Storage
- Container stays alive

Morning (8:30am):
- Girlfriend opens app → Likely hits same container
- Sees your new recipe immediately ✅

Afternoon (2pm):
- Container timed out after 15 min idle
- You open app → New container starts
- Downloads latest DB with all morning changes ✅

Evening (6pm):
- Same container still running
- Both of you use same container
- All changes visible immediately ✅
```

### Unlikely Edge Case (Multiple Containers)

```
- Both of you happen to open app at exact same time
- Cloud Run spawns 2 containers (A and B)
- You add recipe in Container A (uploads to Cloud Storage)
- Girlfriend's next request hits Container B
- Container B hasn't refreshed its /tmp/recipes.db
- She sees old data for ~15 minutes
- After timeout, new container starts with fresh data
```

**Likelihood:** Extremely rare with 2 users making occasional requests.

---

## What Happens On Write

```go
func CreateRecipe(ctx context.Context, recipe *Recipe) error {
    dbMutex.Lock()
    defer dbMutex.Unlock()

    // 1. Write to local /tmp/recipes.db immediately
    result, err := db.ExecContext(ctx, "INSERT INTO recipes ...")

    // 2. Return recipe to user (fast response)
    recipe.ID = id

    // 3. Upload to Cloud Storage asynchronously (doesn't block)
    go func() {
        uploadDBToGCS(context.Background())
        log.Println("Database uploaded to Cloud Storage")
    }()

    return nil
}
```

**Timeline:**
- Write to SQLite: ~1-5ms
- Response to user: ~10-50ms total
- Upload to Cloud Storage: ~100-500ms (happens in background)

**Note:** Other containers don't see this change until they restart and re-download.

---

## Why This Pattern Works for You

### ✅ Perfect For Your Use Case Because:

1. **Small user base (2 users)**
   - Low traffic = likely one container most of the time
   - Cloud Run scales to zero after idle = fresh container on restart

2. **Infrequent writes**
   - You're not adding recipes every second
   - Probably a few recipes per week at most
   - Plenty of time for containers to sync

3. **Casual access pattern**
   - Not a high-stakes transactional system
   - Seeing a recipe 5 minutes late is fine
   - Refresh page = new container = latest data

4. **Ultra-low cost**
   - Nearly free (~$0.001/month)
   - No database server to pay for
   - Cloud Run free tier handles your load

### ❌ When This Pattern Would NOT Work:

- High concurrent writes (>10 users writing simultaneously)
- Need for strong consistency (banking, reservations, etc.)
- Real-time collaboration (multiple users editing same document)
- Large database (>100MB would be slow to sync)
- Mission-critical data where staleness is unacceptable

---

## Recommended Implementation

### Step 1: Add Single Instance Constraint

When you add Cloud Run service to Terraform (`infra/main.tf`), include:

```hcl
resource "google_cloud_run_service" "backend" {
  name     = "recipebook-backend"
  location = var.region

  template {
    spec {
      service_account_name = google_service_account.backend.email

      containers {
        image = "gcr.io/${var.project_id}/recipebook-backend"

        env {
          name  = "DB_BUCKET_NAME"
          value = google_storage_bucket.database.name
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale" = "0"  # Scale to zero
        "autoscaling.knative.dev/maxScale" = "1"  # Single instance
      }
    }
  }
}
```

### Step 2: Deploy and Forget

That's it. With `maxScale = 1`:
- No consistency issues
- Still scales to zero (free when idle)
- Perfect for 2 users
- Simple and reliable

---

## Monitoring

### Check Container Behavior

```bash
# View Cloud Run logs
gcloud run services logs read recipebook-backend --region us-central1

# Look for:
# - "Downloaded database from Cloud Storage" (container startup)
# - "Uploaded database to Cloud Storage" (after writes)
```

### Check Database Freshness

```bash
# See when database was last updated
gsutil ls -l gs://your-project-recipebook-db/recipes.db

# Download and inspect
gsutil cp gs://your-project-recipebook-db/recipes.db /tmp/check.db
sqlite3 /tmp/check.db "SELECT COUNT(*) FROM recipes;"
```

### Container Count

Check Cloud Run metrics in GCP Console:
- Go to Cloud Run → recipebook-backend → Metrics
- Look at "Active instances" graph
- With maxScale=1, should never exceed 1

---

## Summary

**Current Behavior:**
- Each container downloads DB once on startup
- Writes upload to Cloud Storage but don't trigger re-downloads
- Multiple containers can have slightly stale data

**Recommended Solution:**
- Set `maxScale = 1` in Cloud Run config
- Forces single container = no consistency issues
- Still free, still scales to zero
- Perfect for 2 users

**No code changes needed** - just add the Terraform annotation when deploying Cloud Run service.

---

## Questions to Consider

- **Do you want to implement maxScale=1 now or later?**
  - Recommend: Do it now when adding Cloud Run to Terraform

- **What if you need more scale later?**
  - You can always remove the maxScale constraint
  - Or migrate to Cloud SQL when you outgrow SQLite

- **How often will containers restart?**
  - After 15 min idle (scales to zero)
  - On new deployments
  - When Cloud Run updates infrastructure

With maxScale=1, you get the simplicity of SQLite with none of the multi-container headaches.
