# Public Read Model

## Overview

This recipe book uses a **public read, authenticated write** model:

- ✅ **Anyone** can view all recipes (no login required)
- 🔒 **Only you and your girlfriend** can create/edit/delete recipes (login required)

This allows you to share your favorite recipes with friends, family, and the internet while maintaining editorial control.

## Why This Model?

### Use Cases
- Share recipe links with friends via text/email
- Let family browse without creating accounts
- Build a public recipe collection
- SEO-friendly (search engines can index recipes)
- Social sharing (recipes show up in link previews)

### Benefits
✅ **Easy sharing** - Just send someone a link
✅ **No barriers** - Readers don't need accounts
✅ **Low traffic** - You expect "a couple people a month"
✅ **Still secure** - Only editors can modify
✅ **Cost effective** - No auth overhead for reads

## Security Model

### What's Public
```
GET /recipes              → ✅ Public
GET /recipes/{id}         → ✅ Public
GET /recipes/search?q=... → ✅ Public
```

Anyone can:
- Browse all recipes
- View recipe details
- Search for ingredients

### What's Protected
```
POST   /recipes     → 🔒 Auth required
PUT    /recipes/{id} → 🔒 Auth required
DELETE /recipes/{id} → 🔒 Auth required
```

Only authenticated users can:
- Create new recipes
- Update existing recipes
- Delete recipes

## Authentication Flow

### For Readers (No Auth)
```
1. Visitor opens your recipe site
2. Frontend fetches recipes from API (no token needed)
3. Recipes display immediately
```

### For Editors (You + Girlfriend)
```
1. Click "Login" button
2. Sign in with Firebase (email/password)
3. Frontend gets Firebase ID token
4. Token sent with create/edit/delete requests
5. Backend verifies token before allowing changes
```

## Example Usage

### Public Reader
```bash
# Anyone can do this (no token needed)
curl https://recipebook-backend-xxx.run.app/recipes
curl https://recipebook-backend-xxx.run.app/recipes/1
curl https://recipebook-backend-xxx.run.app/recipes/search?q=chicken
```

### Authenticated Editor
```bash
# Must have Firebase token
curl -X POST \
  -H "Authorization: Bearer <firebase-token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"New Recipe","type":"food",...}' \
  https://recipebook-backend-xxx.run.app/recipes
```

## Frontend Implications

### Landing Page (No Login)
```jsx
// No auth check needed
function RecipeList() {
  const [recipes, setRecipes] = useState([]);

  useEffect(() => {
    // Direct API call - no token
    fetch('https://api.example.com/recipes')
      .then(res => res.json())
      .then(setRecipes);
  }, []);

  return <div>{recipes.map(...)}</div>;
}
```

### Edit Page (Login Required)
```jsx
function EditRecipe() {
  const user = useAuth(); // Firebase hook

  if (!user) {
    return <div>Please log in to edit recipes</div>;
  }

  const handleSave = async () => {
    const token = await user.getIdToken();

    await fetch('https://api.example.com/recipes', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(recipe),
    });
  };
}
```

## Cost Implications

### Higher Read Traffic
Since reads are public, you might get more traffic than 2 users:
- Friends browsing recipes: 10-20 people/month
- Search engine bots: ?
- Direct links shared: ?

**Expected monthly requests:**
- Your usage: ~1k requests
- Public readers: ~5k requests (couple people/month)
- **Total: ~6k requests**

**Still well within free tier:**
- Cloud Run free tier: 2M requests/month
- You're using: 0.3% of free tier
- **Cost: $0**

### If You Go Viral
If a recipe somehow gets popular:
- 100k requests/month = still free
- 5M requests/month = $0.10 (3M over free tier × $0.40/M)

**You're safe. Don't worry about cost.**

## Security Considerations

### ✅ What's Protected
- Only 2 people can edit (you + girlfriend)
- Firebase Auth tokens are cryptographically secure
- Can't forge tokens to create/edit recipes
- Database versioning protects against mistakes

### ⚠️ Theoretical Concerns
1. **Recipe scraping**
   - Someone could download all your recipes
   - **But that's okay** - you're sharing publicly anyway

2. **DDoS on read endpoints**
   - Someone could spam GET requests
   - **But** 2M free requests/month is hard to exhaust
   - Would need 66k requests/day sustained

3. **Search spam**
   - Someone could spam search queries
   - **But** SQLite full-text search is fast
   - No database writes = low cost

### 🚨 Real Risk: Near Zero
- No sensitive data (just recipes)
- No PII (no user accounts for readers)
- No payment info
- Worst case: Someone reads all your recipes (which you're sharing anyway)

## Monitoring

Check monthly request counts:

```bash
# View Cloud Run metrics
gcloud run services describe recipebook-backend \
  --region us-central1 \
  --format="value(status.url)"

# Check request count
gcloud monitoring time-series list \
  --filter='metric.type="run.googleapis.com/request_count"' \
  --project your-project-id
```

**Alert if:**
- Requests exceed 100k/month (unusual for "couple people")
- High rate of failed auth attempts (someone trying to hack in)
- Unexpected traffic patterns

## Future Enhancements

### Add Private Recipes
If you want some recipes private:

```go
type Recipe struct {
    // ...
    IsPublic bool `json:"isPublic"`
}

// Only show public recipes to unauthenticated users
func GetRecipes(ctx context.Context, authenticated bool) ([]Recipe, error) {
    query := "SELECT * FROM recipes"
    if !authenticated {
        query += " WHERE is_public = true"
    }
    // ...
}
```

### Add Recipe Authors
Track who created each recipe:

```go
type Recipe struct {
    // ...
    CreatedBy string `json:"createdBy"`
}
```

### Add Comments
Let readers comment (with moderation):
- Store comments in separate table
- Require approval before showing
- Or use Disqus/other service

## Comparison to Other Models

### Private (All Auth Required)
```
Pros: Maximum security
Cons: Friends need accounts, can't share links
```

### Public Write (Wiki-style)
```
Pros: Community contributions
Cons: Spam, vandalism, need moderation
```

### **Public Read, Private Write (Current)**
```
Pros: Easy sharing, editorial control, low cost
Cons: Can't accept community recipes (could add this later)
```

## Summary

Your setup is perfect for:
- ✅ Sharing recipes with friends/family
- ✅ Low traffic expectations
- ✅ Maintaining editorial control
- ✅ Zero cost at your scale
- ✅ Simple user experience

No changes needed. The `allUsers` IAM permission makes sense now because you **want** the API to be publicly accessible for reads.

Just make sure only you and your girlfriend have Firebase Auth accounts for editing!
