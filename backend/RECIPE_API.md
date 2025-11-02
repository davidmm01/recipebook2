# Recipe API Documentation

## Authentication Model

**Public Read, Authenticated Write:**
- ✅ Anyone can **read** recipes (GET requests)
- 🔒 Only authenticated users can **create/update/delete** recipes (POST/PUT/DELETE)

This allows you to share your recipes publicly while maintaining control over who can edit them.

## Recipe Format

All recipes use markdown format for text fields.

```json
{
  "id": 1,
  "title": "Chicken Pasta",
  "type": "food",
  "ingredients": "## Ingredients\n\n- 2 chicken breasts\n- 500g pasta\n- 400ml cream\n- Salt and pepper",
  "method": "## Method\n\n1. Cook pasta according to package instructions\n2. Pan-fry chicken until cooked through\n3. Add cream and simmer\n4. Season to taste",
  "notes": "## Notes\n\nGreat for meal prep. Can freeze for up to 3 months.",
  "createdAt": "2025-01-24T12:00:00Z",
  "updatedAt": "2025-01-24T12:00:00Z"
}
```

## Endpoints

### GET /recipes
List all recipes. **Public endpoint - no authentication required.**

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "title": "Chicken Pasta",
    "type": "food",
    "ingredients": "...",
    "method": "...",
    "notes": "...",
    "createdAt": "2025-01-24T12:00:00Z",
    "updatedAt": "2025-01-24T12:00:00Z"
  }
]
```

---

### POST /recipes
Create a new recipe.

**Headers:**
```
Authorization: Bearer <firebase-id-token>
Content-Type: application/json
```

**Body:**
```json
{
  "title": "Chicken Pasta",
  "type": "food",
  "ingredients": "## Ingredients\n\n- 2 chicken breasts\n- 500g pasta",
  "method": "## Method\n\n1. Cook pasta...",
  "notes": "Optional notes"
}
```

**Response:** `201 Created`
```json
{
  "id": 1,
  "title": "Chicken Pasta",
  "type": "food",
  "ingredients": "...",
  "method": "...",
  "notes": "...",
  "createdAt": "2025-01-24T12:00:00Z",
  "updatedAt": "2025-01-24T12:00:00Z"
}
```

---

### GET /recipes/{id}
Get a single recipe by ID. **Public endpoint - no authentication required.**

**Response:** `200 OK`
```json
{
  "id": 1,
  "title": "Chicken Pasta",
  "type": "food",
  "ingredients": "...",
  "method": "...",
  "notes": "...",
  "createdAt": "2025-01-24T12:00:00Z",
  "updatedAt": "2025-01-24T12:00:00Z"
}
```

**Error:** `404 Not Found` if recipe doesn't exist

---

### PUT /recipes/{id}
Update an existing recipe.

**Headers:**
```
Authorization: Bearer <firebase-id-token>
Content-Type: application/json
```

**Body:**
```json
{
  "title": "Updated Title",
  "type": "food",
  "ingredients": "...",
  "method": "...",
  "notes": "..."
}
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "title": "Updated Title",
  "type": "food",
  "ingredients": "...",
  "method": "...",
  "notes": "...",
  "createdAt": "2025-01-24T12:00:00Z",
  "updatedAt": "2025-01-24T14:30:00Z"
}
```

**Error:** `404 Not Found` if recipe doesn't exist

---

### DELETE /recipes/{id}
Delete a recipe.

**Headers:**
```
Authorization: Bearer <firebase-id-token>
```

**Response:** `204 No Content`

**Error:** `404 Not Found` if recipe doesn't exist

---

### GET /recipes/search?q=query
Full-text search across all recipe fields. **Public endpoint - no authentication required.**

**Query Parameters:**
- `q` (required): Search query

**Examples:**
- `/recipes/search?q=chicken` - Find recipes containing "chicken"
- `/recipes/search?q=pasta AND tomato` - Find recipes with both terms
- `/recipes/search?q="chicken breast"` - Exact phrase search

**Response:** `200 OK`
```json
[
  {
    "id": 1,
    "title": "Chicken Pasta",
    "type": "food",
    "ingredients": "...",
    "method": "...",
    "notes": "...",
    "createdAt": "2025-01-24T12:00:00Z",
    "updatedAt": "2025-01-24T12:00:00Z"
  }
]
```

Results are ranked by relevance.

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

- All recipes are shared across authenticated users
- Firebase Auth is required for all endpoints except `/health`
- Markdown formatting is supported in `ingredients`, `method`, and `notes` fields
- Full-text search uses SQLite FTS5 for fast, relevant results
- Recipe updates modify the `updatedAt` timestamp automatically
