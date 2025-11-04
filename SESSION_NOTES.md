# Session Notes - November 3, 2025

## What We Accomplished

### Recipe Import Tool ✅
Created a complete system to import YAML recipes into the SQLite database:

1. **Import Script**: `backend/cmd/import-recipes/main.go`
   - Converts YAML files from `backend/data/cocktails/` and `backend/data/recipes/`
   - Generates SQL INSERT statements with proper mappings
   - Successfully processes all 38 recipes (4 cocktails, 34 food recipes)

2. **Special Features Working**:
   - ✅ YAML comments (e.g., `# sauces`) converted to markdown subheadings (`### sauces`)
   - ✅ Notes and Next fields combined with proper subheadings
   - ✅ Tags properly linked through many-to-many relationship
   - ✅ Recipe type auto-detected from directory (food/drink)
   - ✅ All source fields combined (except submitter)
   - ✅ Created dates and creator names preserved

3. **Makefile Targets Added**:
   ```bash
   make import-recipes      # Generate import.sql
   make import-and-load     # Generate and import in one step
   ```

4. **Documentation**: Created `backend/data/README.md` with full instructions

### Example Output
**Ez Mapo Tofu** - Shows ingredient subheadings working:
```markdown
- neutral oil
- 1 brown onion, sliced finely
...
### sauces
- 113g of Lee Kum Kee mapo tofu sauce
- 1 tablespoon oyster sauce
```

**Healthy Curry** - Shows multiple subheadings:
```markdown
### Protein & marinade
...
### Base
...
### Spices for oil
...
```

## Current Issue 🔍

### User Role Lost After Logout/Login

**Symptom**: After logging out and back in, your role was forgotten

**Investigation Results**:
- Database shows you still have "admin" role (firebase_uid: susKfNjTddc8x4KbHCHwIekb05e2)
- Database query: `sqlite3 /tmp/recipes.db "SELECT firebase_uid, email, role FROM users;"`
  Result: `susKfNjTddc8x4KbHCHwIekb05e2|d.m.michieli@gmail.com|admin`
- Backend is working correctly and has your role stored

**Likely Cause**:
The issue appears to be on the **frontend side**. The backend code in `main.go:467` creates new users with "viewer" role by default, but this only happens if the user doesn't exist in the database. Since your user exists with admin role, the backend should be returning it correctly.

**Next Steps**:
1. Check frontend authentication/user context code
2. Look for how the role is cached/stored in browser (localStorage, state, etc.)
3. Verify the `/user/profile` API endpoint is being called after login
4. Check browser console for any errors during login

## Important Notes

### Database Import Workflow
**Critical**: The backend must be running BEFORE importing recipes!

1. Start backend: `make run-local`
2. Wait for "Database initialized successfully" in logs
3. Import recipes: `make import-and-load`

**Why**: The backend downloads/creates the database on startup. If you import before starting the backend, the data gets overwritten.

## Files Modified/Created

### New Files
- `backend/cmd/import-recipes/main.go` - Import script
- `backend/data/README.md` - Import documentation
- `backend/import.sql` - Generated SQL (can be regenerated anytime)

### Modified Files
- `backend/Makefile` - Added `import-recipes` and `import-and-load` targets
- `backend/go.mod` - Added `gopkg.in/yaml.v3` dependency

## Quick Commands for Next Session

```bash
# Check your user role in database
sqlite3 /tmp/recipes.db "SELECT firebase_uid, email, role FROM users WHERE email='d.m.michieli@gmail.com';"

# Re-import recipes (if needed)
make import-and-load

# Check recipe count
sqlite3 /tmp/recipes.db "SELECT COUNT(*) FROM recipes;"

# Test API
curl http://localhost:8080/recipes | python3 -m json.tool | head -50
```

## Backend Status
- Backend is currently running (PID: 11334)
- Port 8080 is in use
- Database at `/tmp/recipes.db` contains 38 recipes with all your data intact
