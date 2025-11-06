# Recipe API Documentation

## Authentication Model

**Public Read, Authenticated Write:**
- ✅ Anyone can **read** recipes (GET requests)
- 🔒 Only authenticated users can **create/update/delete** recipes (POST/PUT/DELETE)

This allows you to share your recipes publicly while maintaining control over who can edit them.

## Recipe Format

All recipes use markdown format for text fields and UUID strings for IDs.

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Chicken Pasta",
  "description": "Creamy pasta with pan-fried chicken breast",
  "type": "food",
  "cuisine": "italian",
  "tags": ["pasta", "quick", "italian"],
  "ingredients": "## Ingredients\n\n- 2 chicken breasts\n- 500g pasta\n- 400ml cream\n- Salt and pepper",
  "method": "## Method\n\n1. Cook pasta according to package instructions\n2. Pan-fry chicken until cooked through\n3. Add cream and simmer\n4. Season to taste",
  "notes": "## Notes\n\nGreat for meal prep. Can freeze for up to 3 months.",
  "images": [
    {
      "id": 1,
      "recipeId": "550e8400-e29b-41d4-a716-446655440000",
      "imageUrl": "https://storage.googleapis.com/bucket/recipe-images/550e8400.../image.jpg",
      "displayOrder": 0,
      "createdAt": "2025-01-24T12:00:00Z"
    }
  ],
  "createdByUserId": "firebase-uid-abc123",
  "createdByName": "John Doe",
  "createdAt": "2025-01-24T12:00:00Z",
  "updatedAt": "2025-01-24T12:00:00Z"
}
```

**New Creator Fields:**
- `createdByUserId` (string, nullable): Firebase UID of the user who created the recipe
- `createdByName` (string, nullable): Display name of the user who created the recipe
- Both fields are optional and will be `null` for legacy recipes or anonymous submissions

## Endpoints

### GET /recipes
List all recipes. **Public endpoint - no authentication required.**

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Chicken Pasta",
    "description": "Creamy pasta with pan-fried chicken breast",
    "type": "food",
    "cuisine": "italian",
    "tags": ["pasta", "quick", "italian"],
    "ingredients": "...",
    "method": "...",
    "notes": "...",
    "images": [
      {
        "id": 1,
        "recipeId": "550e8400-e29b-41d4-a716-446655440000",
        "imageUrl": "https://storage.googleapis.com/bucket/recipe-images/550e8400.../image.jpg",
        "displayOrder": 0,
        "createdAt": "2025-01-24T12:00:00Z"
      }
    ],
    "createdAt": "2025-01-24T12:00:00Z",
    "updatedAt": "2025-01-24T12:00:00Z"
  }
]
```

---

### POST /recipes
Create a new recipe. **Requires authentication (editor or admin role).**

**Headers:**
```
Authorization: Bearer <firebase-id-token>
Content-Type: application/json
```

**Body:**
```json
{
  "title": "Chicken Pasta",
  "description": "Creamy pasta with pan-fried chicken breast",
  "type": "food",
  "cuisine": "italian",
  "tags": ["pasta", "quick", "italian"],
  "ingredients": "## Ingredients\n\n- 2 chicken breasts\n- 500g pasta",
  "method": "## Method\n\n1. Cook pasta...",
  "notes": "Optional notes"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Chicken Pasta",
  "description": "Creamy pasta with pan-fried chicken breast",
  "type": "food",
  "cuisine": "italian",
  "tags": ["pasta", "quick", "italian"],
  "ingredients": "...",
  "method": "...",
  "notes": "...",
  "images": [],
  "createdAt": "2025-01-24T12:00:00Z",
  "updatedAt": "2025-01-24T12:00:00Z"
}
```

---

### GET /recipes/{id}
Get a single recipe by UUID. **Public endpoint - no authentication required.**

**Example:** `GET /recipes/550e8400-e29b-41d4-a716-446655440000`

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Chicken Pasta",
  "description": "Creamy pasta with pan-fried chicken breast",
  "type": "food",
  "cuisine": "italian",
  "tags": ["pasta", "quick", "italian"],
  "ingredients": "...",
  "method": "...",
  "notes": "...",
  "images": [
    {
      "id": 1,
      "recipeId": "550e8400-e29b-41d4-a716-446655440000",
      "imageUrl": "https://storage.googleapis.com/bucket/recipe-images/550e8400.../image.jpg",
      "displayOrder": 0,
      "createdAt": "2025-01-24T12:00:00Z"
    }
  ],
  "createdAt": "2025-01-24T12:00:00Z",
  "updatedAt": "2025-01-24T12:00:00Z"
}
```

**Error:** `404 Not Found` if recipe doesn't exist

---

### PUT /recipes/{id}
Update an existing recipe. **Requires authentication (editor or admin role).**

**Headers:**
```
Authorization: Bearer <firebase-id-token>
Content-Type: application/json
```

**Body:**
```json
{
  "title": "Updated Title",
  "description": "Updated description",
  "type": "food",
  "cuisine": "italian",
  "tags": ["pasta", "updated"],
  "ingredients": "...",
  "method": "...",
  "notes": "..."
}
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Updated Title",
  "description": "Updated description",
  "type": "food",
  "cuisine": "italian",
  "tags": ["pasta", "updated"],
  "ingredients": "...",
  "method": "...",
  "notes": "...",
  "images": [],
  "createdAt": "2025-01-24T12:00:00Z",
  "updatedAt": "2025-01-24T14:30:00Z"
}
```

**Error:** `404 Not Found` if recipe doesn't exist

---

### DELETE /recipes/{id}
Delete a recipe. **Requires authentication (admin role).**

**Headers:**
```
Authorization: Bearer <firebase-id-token>
```

**Response:** `204 No Content`

**Error:** `404 Not Found` if recipe doesn't exist

---

### GET /recipes/search?q=query
Full-text search across all recipe text fields. **Public endpoint - no authentication required.**

**Query Parameters:**
- `q` (required): Search query

**Searchable Fields:**
- title
- description
- cuisine
- ingredients
- method
- notes

**Examples:**
- `/recipes/search?q=chicken` - Find recipes containing "chicken"
- `/recipes/search?q=pasta AND tomato` - Find recipes with both terms
- `/recipes/search?q="chicken breast"` - Exact phrase search
- `/recipes/search?q=italian` - Find Italian recipes (searches cuisine field)

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Chicken Pasta",
    "description": "Creamy pasta with pan-fried chicken breast",
    "type": "food",
    "cuisine": "italian",
    "tags": ["pasta", "quick", "italian"],
    "ingredients": "...",
    "method": "...",
    "notes": "...",
    "images": [],
    "createdAt": "2025-01-24T12:00:00Z",
    "updatedAt": "2025-01-24T12:00:00Z"
  }
]
```

Results are ranked by relevance using SQLite FTS5.

---

### POST /recipes/images
Upload an image for a recipe. **Requires authentication (editor or admin role).**

**Headers:**
```
Authorization: Bearer <firebase-id-token>
Content-Type: multipart/form-data
```

**Form Fields:**
- `recipeId` (required): UUID of the recipe to attach the image to
- `image` (required): Image file (max 10MB)
- `displayOrder` (optional): Integer for ordering multiple images (default: 0)

**Supported Formats:**
- `.jpg` / `.jpeg`
- `.png`
- `.gif`
- `.webp`

**Response:** `200 OK`
```json
{
  "imageUrl": "https://storage.googleapis.com/bucket/recipe-images/550e8400.../abc123.jpg"
}
```

**Notes:**
- Images are stored in Google Cloud Storage
- Each recipe can have multiple images
- Images are publicly accessible via the returned URL
- Image URLs are automatically included in recipe responses
- If database save fails, the uploaded image is automatically cleaned up from storage

**Errors:**
- `400 Bad Request` - Invalid file type, missing fields, or file too large
- `401 Unauthorized` - Missing or invalid authentication token
- `500 Internal Server Error` - Upload or storage failure

---

## User Profile Endpoints

### GET /user/profile
Get the authenticated user's profile. **Requires authentication.**

**Headers:**
```
Authorization: Bearer <firebase-id-token>
```

**Response:** `200 OK`
```json
{
  "email": "user@example.com",
  "displayName": "John Doe",
  "role": "editor",
  "lastLoginAt": "2025-01-24T12:00:00Z"
}
```

**Notes:**
- Returns the profile of the currently authenticated user
- `displayName` can be updated via PUT request
- User profiles are stored in SQLite database
- Users are automatically created on first login with default "viewer" role

---

### PUT /user/profile
Update the authenticated user's display name. **Requires authentication.**

**Headers:**
```
Authorization: Bearer <firebase-id-token>
Content-Type: application/json
```

**Body:**
```json
{
  "displayName": "Jane Smith"
}
```

**Response:** `200 OK`
```json
{
  "email": "user@example.com",
  "displayName": "Jane Smith",
  "role": "editor",
  "lastLoginAt": "2025-01-24T12:00:00Z"
}
```

**Notes:**
- Only `displayName` can be updated by the user via this endpoint
- Role changes must be done via CLI tool (see `backend/cmd/manage-users/`)
- Email cannot be changed (tied to Firebase Auth account)
- Display name is automatically used when creating recipes

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "Missing authorization header"
}
```

### 404 Not Found
```json
{
  "error": "Recipe not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to get recipes"
}
```

---

## Notes

### Access Control
- All recipes are **public for reading** (GET requests require no authentication)
- **Write operations** require authentication and appropriate role:
  - **Viewers**: Can only read recipes
  - **Editors**: Can create and update recipes
  - **Admins**: Can delete recipes and manage user roles
- User profiles and roles are stored in SQLite database
- **Authentication**: Firebase Authentication for Google login
- **Authorization**: Role-based access control via SQLite `users` table
- Users are auto-created on first login with "viewer" role

### Data Format
- Recipe IDs are **UUID strings**, not integers
- Markdown formatting is supported in `ingredients`, `method`, `notes`, and `description` fields
- Tags are stored as lowercase strings for consistency
- Multiple tags can be assigned to each recipe

### Search
- Full-text search uses SQLite **FTS5** for fast, relevant results
- Searches across: title, description, cuisine, ingredients, method, and notes
- Results are ranked by relevance

### Images
- Images are stored in **Google Cloud Storage**, not in the database
- Each recipe can have **multiple images** with configurable display order
- Images are **publicly accessible** via GCS URLs
- Maximum file size: **10MB**
- Supported formats: `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`

### Timestamps
- Recipe creation sets `createdAt`
- Recipe updates automatically modify `updatedAt`

### Creator Tracking
- New recipes automatically capture the creator's **Firebase UID** in `createdByUserId`
- If the user has set a **display name**, it's captured in `createdByName`
- Both fields are **nullable** to support:
  - Legacy recipes created before this feature
  - Anonymous recipe submissions (if enabled)
  - Recipes without attribution
- User display names can be managed via `/user/profile` endpoint

### User Role Management
- User roles are managed via **CLI tool only** (not through the API)
- CLI tool location: `backend/cmd/manage-users/`
- Commands:
  - `./manage-users list` - List all users
  - `./manage-users set-role --email <email> --role <role>` - Change user role
- See `backend/cmd/manage-users/README.md` for full documentation
- This prevents users from escalating their own privileges
