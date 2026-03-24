# Plan: Injectable Auth for Testable Handlers

## Context

`authenticateRequest` is called directly in every protected handler. Because it makes a live
network call to Firebase to verify tokens, tests can never get past the 401 — making it
impossible to test any logic that runs *after* authentication (role checks, ownership checks,
happy paths). The fix is to put `authenticateRequest` behind an injectable function variable
so tests can swap it out for a stub.

---

## Approach

Introduce a package-level `authenticate` variable of a named function type. All handlers call
`authenticate(r)` instead of `authenticateRequest(r)` directly. Tests override `authenticate`
to return a predetermined UID without touching Firebase.

### Step 1 — Declare the type and variable (`backend/main.go`)

Add near the top of `main.go`, alongside other package-level vars:

```go
// AuthFunc validates a request and returns the caller's Firebase UID.
type AuthFunc func(r *http.Request) (string, error)

// authenticate is the auth function used by all handlers.
// Overridden in tests to avoid live Firebase calls.
var authenticate AuthFunc = authenticateRequest
```

### Step 2 — Replace all 11 call sites

Replace every `authenticateRequest(r)` with `authenticate(r)`.

Call sites (all currently `authenticateRequest`):
| File | Handler | Line (approx) |
|------|---------|---------------|
| `backend/main.go` | `recipesHandler` POST | ~108 |
| `backend/main.go` | `recipeByIDHandler` PUT | ~170 |
| `backend/main.go` | `recipeByIDHandler` DELETE | ~195 |
| `backend/main.go` | `imageUploadHandler` | ~312 |
| `backend/main.go` | `iconsHandler` POST | ~393 |
| `backend/main.go` | `makeLogsHandler` POST | ~473 |
| `backend/main.go` | `makeLogByIDHandler` PUT | ~517 |
| `backend/main.go` | `makeLogByIDHandler` DELETE | ~543 |
| `backend/admin.go` | `requireAdmin` | ~13 |
| `backend/user.go` | `userProfileHandler` | ~15 |

### Step 3 — Test helpers (`backend/main_test.go` or new `backend/auth_test_helpers_test.go`)

Add a helper to set up a fake auth context for a given UID, and a cleanup function:

```go
func withAuth(uid string) func() {
    original := authenticate
    authenticate = func(r *http.Request) (string, error) {
        return uid, nil
    }
    return func() { authenticate = original }
}
```

Usage in tests:
```go
defer withAuth("some-uid")()
// now call handler — it will see "some-uid" as the authenticated user
```

### Step 4 — New tests enabled by this change

**`backend/main_test.go`** — Recipe delete permission tests:
- Admin can delete any recipe (seed recipe with different creator UID, auth as admin)
- Owner can delete their own recipe (seed recipe with matching creator UID)
- Non-owner non-admin gets 403 (seed recipe with different creator, auth as viewer)

**`backend/main_test.go`** — Admin handler tests:
- Non-admin calling `GET /admin/users` gets 403 (auth as viewer/editor)
- Admin calling `GET /admin/users` gets 200 with user list
- Admin calling `PUT /admin/users/{uid}/role` succeeds
- Admin calling `PUT /admin/users/{own-uid}/role` gets 403 (self-change prevention)

---

## Files Modified

| File | Change |
|------|--------|
| `backend/main.go` | Add `AuthFunc` type + `authenticate` var; replace 8 call sites |
| `backend/admin.go` | Replace 1 call site in `requireAdmin` |
| `backend/user.go` | Replace 1 call site in `userProfileHandler` |
| `backend/main_test.go` | Add `withAuth` helper + new tests for delete perms and admin handlers |

---

## Verification

```bash
cd backend && go test -tags fts5 ./...
```

All existing tests should still pass. New tests should cover the previously unreachable
post-auth logic paths.
