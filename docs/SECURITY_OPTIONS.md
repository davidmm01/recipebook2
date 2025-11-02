# Cloud Run Security Options

This document explains the security model and alternative configurations.

## Current Setup: Public Endpoint + Firebase Auth

### Configuration
```hcl
# Cloud Run is publicly accessible
resource "google_cloud_run_v2_service_iam_member" "public_access" {
  role   = "roles/run.invoker"
  member = "allUsers"
}
```

### How It Works
```
Internet → Cloud Run (public) → Firebase Auth Check → Protected Data
```

**Anyone can reach the endpoint, but:**
- All `/recipes` endpoints require valid Firebase ID token
- Only `/health` is unprotected

### Pros
✅ Simple - frontend just calls API directly
✅ No additional auth complexity
✅ Standard pattern for Firebase Auth apps
✅ Easy local development

### Cons
❌ API endpoint is discoverable
❌ No built-in DDoS protection
❌ Health endpoint can be spammed
❌ Could enumerate endpoints

---

## Option 1: Remove Public Access (More Secure)

### Configuration
```hcl
# Remove this resource entirely:
# resource "google_cloud_run_v2_service_iam_member" "public_access" { ... }

# Keep the service private by default
```

### How It Works
```
Internet → Cloud Run rejects (401) ✅

Authenticated request → Cloud Run (checks service account) → Firebase Auth → Data
```

**Only requests with valid Google Cloud credentials can reach the service.**

### Implementation

**Frontend must get ID token from Cloud Run:**

```javascript
// In frontend
async function callAPI() {
  // Get ID token for Cloud Run service
  const targetAudience = 'https://recipebook-backend-xxx.run.app';

  // This requires the user to be signed in to Google Cloud
  // NOT the same as Firebase Auth!
  const idToken = await getCloudRunIdToken(targetAudience);

  const response = await fetch(`${API_URL}/recipes`, {
    headers: {
      'Authorization': `Bearer ${idToken}`,
    }
  });
}
```

### Pros
✅ Cloud Run endpoint not publicly accessible
✅ Better DDoS protection
✅ Can't enumerate endpoints

### Cons
❌ Much more complex - need TWO auth systems
❌ Users must sign in with Google Cloud (not just Firebase)
❌ Harder to develop locally
❌ Overkill for your use case

---

## Option 2: Firebase Auth + Cloud Armor (Recommended for Production)

### Configuration

Keep public access, but add Cloud Armor for protection:

```hcl
# Keep public access
resource "google_cloud_run_v2_service_iam_member" "public_access" {
  member = "allUsers"
}

# Add Cloud Armor security policy
resource "google_compute_security_policy" "policy" {
  name = "recipebook-security-policy"

  # Rate limiting
  rule {
    action   = "rate_based_ban"
    priority = 1000
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    rate_limit_options {
      conform_action = "allow"
      exceed_action  = "deny(429)"
      enforce_on_key = "IP"
      rate_limit_threshold {
        count        = 100
        interval_sec = 60
      }
      ban_duration_sec = 600
    }
  }

  # Block known bad IPs
  rule {
    action   = "deny(403)"
    priority = 2000
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["0.0.0.0/0"]  # Add specific bad IPs here
      }
    }
  }

  # Default allow
  rule {
    action   = "allow"
    priority = 2147483647
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
  }
}
```

### Pros
✅ Rate limiting prevents abuse
✅ Can block malicious IPs
✅ Still simple for frontend
✅ Production-grade security

### Cons
❌ Cloud Armor costs money (~$1-2/month minimum)
❌ More complex Terraform
❌ Might be overkill for 2 users

---

## Option 3: Firebase App Check (Best Balance)

### Configuration

Add App Check to verify requests come from your app:

```javascript
// In frontend
import { initializeAppCheck, ReCaptchaV3Provider } from 'firebase/app-check';

const appCheck = initializeAppCheck(app, {
  provider: new ReCaptchaV3Provider('your-recaptcha-site-key'),
  isTokenAutoRefreshEnabled: true
});
```

```go
// In backend
import (
  "firebase.google.com/go/v4/appcheck"
)

func AppCheckMiddleware(next http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    // Verify App Check token
    appCheckToken := r.Header.Get("X-Firebase-AppCheck")

    _, err := appCheckClient.VerifyToken(appCheckToken)
    if err != nil {
      http.Error(w, "Invalid App Check token", 401)
      return
    }

    next.ServeHTTP(w, r)
  }
}
```

### Pros
✅ Prevents bot/scraper abuse
✅ Verifies requests from your app
✅ Free tier: 10k verifications/month
✅ Integrates with Firebase
✅ Simple to implement

### Cons
❌ Requires reCAPTCHA setup
❌ May annoy users with captchas
❌ Another dependency

---

## Recommendation for Your Use Case

### For Development (Now)
**Keep public access** - it's simple and works

```hcl
# Current setup is fine
resource "google_cloud_run_v2_service_iam_member" "public_access" {
  member = "allUsers"
}
```

**Why:**
- Only 2 users (you + girlfriend)
- Firebase Auth protects all important endpoints
- Unlikely to be targeted by attackers
- Free tier has 2M requests/month (more than enough)

### For Production (Later, if needed)

**Add Firebase App Check** when:
- You get unexpected traffic
- You notice abuse in logs
- You want to prevent bots

**OR add Cloud Armor** when:
- You're getting DDoS'd
- You want rate limiting
- Cost isn't a concern (~$1-2/month)

---

## What's Actually at Risk?

Let's be realistic about threats:

### ✅ **Your Data is Safe**
- All recipes protected by Firebase Auth
- Tokens can't be forged
- User IDs verified

### ⚠️ **Theoretical Concerns**
1. **DDoS on /health endpoint**
   - Could spam requests
   - You have 2M free requests/month
   - Would need 33k requests/day to exceed
   - Unlikely for a personal app

2. **Endpoint enumeration**
   - Attacker could discover `/recipes`, `/recipes/search`, etc.
   - But can't access data without auth
   - Not a real risk

3. **Credential stuffing**
   - Attacker could try leaked passwords against your login
   - Firebase has built-in rate limiting
   - Only 2 user accounts to target
   - Very unlikely

### 🚨 **Real Risk: Zero**
Your recipe data is safe. The worst that can happen:
- Someone spams `/health` and wastes your free tier requests
- Someone discovers your API structure (but can't access data)

---

## How to Change It

### Option A: Make it Private (Remove Public Access)

```hcl
# In infra/main.tf - DELETE this entire resource:
# resource "google_cloud_run_v2_service_iam_member" "public_access" {
#   name     = google_cloud_run_v2_service.backend.name
#   location = google_cloud_run_v2_service.backend.location
#   role     = "roles/run.invoker"
#   member   = "allUsers"
# }
```

Then frontend needs Cloud Run authentication (complex).

### Option B: Add Firebase App Check

1. Enable App Check in Firebase Console
2. Add App Check SDK to frontend
3. Verify App Check tokens in backend
4. Keep public access for Cloud Run

### Option C: Keep Current Setup

Do nothing. Your data is protected by Firebase Auth.

Monitor Cloud Run metrics. If you see abuse, add protections then.

---

## My Recommendation

**Keep the current setup.**

Here's why:
- Your data is already protected (Firebase Auth)
- 2 users won't attract attackers
- Free tier is generous (2M requests/month)
- You can monitor and add protection if needed
- Don't over-engineer before you have a problem

**Monitor for issues:**
```bash
# Check request counts monthly
gcloud monitoring time-series list \
  --filter='metric.type="run.googleapis.com/request_count"' \
  --project your-project-id

# If you consistently exceed 100k requests/month (with 2 users), investigate
```

**Add Firebase App Check later if:**
- You see suspicious traffic patterns
- You exceed 50k requests/month (unusual for 2 users)
- You open it to more users
- You just want extra security

The current setup is **secure for your use case**. Don't worry about it.
