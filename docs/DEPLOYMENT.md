# Deployment Guide

This guide walks you through deploying the RecipeBook application to GCP.

## Prerequisites

1. **GCP Account** with billing enabled
2. **gcloud CLI** installed and authenticated
3. **Terraform** installed (v1.0+)
4. **Firebase Project** created (for authentication)

## Step-by-Step Deployment

### 1. Set up GCP Project

```bash
# Set your project ID
export PROJECT_ID="your-gcp-project-id"

# Set the project
gcloud config set project $PROJECT_ID

# Enable required APIs
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable storage.googleapis.com
```

### 2. Deploy Infrastructure with Terraform

```bash
cd infra

# Create terraform.tfvars
cat > terraform.tfvars <<EOF
project_id = "$PROJECT_ID"
region     = "us-central1"
EOF

# Initialize Terraform
terraform init

# Review the plan
terraform plan

# Apply infrastructure
terraform apply
```

This creates:
- Cloud Storage bucket for SQLite database
- Service account for Cloud Run
- Cloud Run service (with placeholder image)

**Save the outputs:**
```bash
# Get the bucket name
export DB_BUCKET_NAME=$(terraform output -raw database_bucket_name)

# Get the service account email
export SERVICE_ACCOUNT=$(terraform output -raw backend_service_account_email)
```

### 3. Build and Deploy Backend

```bash
cd ../backend

# Build and deploy to Cloud Run
gcloud run deploy recipebook-backend \
  --source . \
  --region us-central1 \
  --service-account $SERVICE_ACCOUNT \
  --allow-unauthenticated \
  --max-instances 1 \
  --min-instances 0
```

The `--source .` flag will:
1. Build the Docker image from the Dockerfile
2. Push it to Google Container Registry
3. Deploy to Cloud Run

**Note:** The `--max-instances 1` is crucial - it ensures only one container runs, avoiding consistency issues.

### 4. Verify Backend is Running

```bash
# Get the Cloud Run URL
export BACKEND_URL=$(gcloud run services describe recipebook-backend \
  --region us-central1 \
  --format 'value(status.url)')

# Test the health endpoint
curl $BACKEND_URL/health

# Should return: {"status":"healthy"}
```

### 5. Set up Firebase Authentication

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Select your project
3. Go to **Authentication** → **Sign-in method**
4. Enable **Email/Password**
5. Go to **Project Settings** → **General**
6. Under "Your apps", add a web app
7. Copy the Firebase config (you'll need this for frontend)

### 6. Download Firebase Service Account (for Backend)

The backend needs a service account to verify Firebase tokens:

1. Firebase Console → **Project Settings** → **Service Accounts**
2. Click **Generate new private key**
3. **DO NOT commit this file to git!**

For Cloud Run, upload it as a secret (optional, for production):

```bash
# Create secret
gcloud secrets create firebase-service-account \
  --data-file=service-account.json

# Grant access to Cloud Run service account
gcloud secrets add-iam-policy-binding firebase-service-account \
  --member="serviceAccount:$SERVICE_ACCOUNT" \
  --role="roles/secretmanager.secretAccessor"
```

Or simpler: Use Application Default Credentials (Cloud Run can access Firebase Admin SDK automatically if using same project).

### 7. Deploy Frontend (TODO)

Frontend deployment to Firebase Hosting will be added later.

For now, you can run locally:

```bash
cd frontend

# Create .env file with Firebase config
cat > .env <<EOF
REACT_APP_API_URL=$BACKEND_URL
REACT_APP_FIREBASE_API_KEY=your-api-key
REACT_APP_FIREBASE_AUTH_DOMAIN=your-project.firebaseapp.com
REACT_APP_FIREBASE_PROJECT_ID=your-project-id
REACT_APP_FIREBASE_STORAGE_BUCKET=your-project.appspot.com
REACT_APP_FIREBASE_MESSAGING_SENDER_ID=your-sender-id
REACT_APP_FIREBASE_APP_ID=your-app-id
EOF

npm install
npm start
```

## Updating the Backend

After making code changes:

```bash
cd backend

# Redeploy
gcloud run deploy recipebook-backend \
  --source . \
  --region us-central1
```

## Monitoring

### View Logs

```bash
# Stream Cloud Run logs
gcloud run services logs tail recipebook-backend --region us-central1

# View in Cloud Console
echo "https://console.cloud.google.com/run/detail/us-central1/recipebook-backend/logs?project=$PROJECT_ID"
```

### Check Database

```bash
# Download current database
gsutil cp gs://$DB_BUCKET_NAME/recipes.db /tmp/recipes.db

# Inspect with SQLite
sqlite3 /tmp/recipes.db "SELECT COUNT(*) FROM recipes;"
sqlite3 /tmp/recipes.db "SELECT * FROM recipes LIMIT 5;"
```

### View Metrics

```bash
# Open Cloud Run metrics
echo "https://console.cloud.google.com/run/detail/us-central1/recipebook-backend/metrics?project=$PROJECT_ID"
```

## Cost Monitoring

```bash
# Check current month billing
gcloud billing accounts list
gcloud billing projects describe $PROJECT_ID

# View in console
echo "https://console.cloud.google.com/billing?project=$PROJECT_ID"
```

Expected costs: **~$0.001/month** for 2 users with occasional usage.

## Troubleshooting

### Backend won't start

```bash
# Check logs for errors
gcloud run services logs read recipebook-backend --region us-central1 --limit 50

# Common issues:
# - Missing DB_BUCKET_NAME environment variable
# - Service account doesn't have access to bucket
# - SQLite CGO build issue
```

### Database not syncing

```bash
# Verify bucket exists
gsutil ls gs://$DB_BUCKET_NAME/

# Check service account permissions
gsutil iam get gs://$DB_BUCKET_NAME/

# Should show recipebook-backend service account with objectAdmin role
```

### Can't authenticate

```bash
# Verify Firebase Auth is enabled in console
# Check that frontend is using correct Firebase config
# Verify backend can reach Firebase Auth API
```

## Rollback

If you need to rollback a deployment:

```bash
# List revisions
gcloud run revisions list --service recipebook-backend --region us-central1

# Rollback to previous revision
gcloud run services update-traffic recipebook-backend \
  --to-revisions REVISION_NAME=100 \
  --region us-central1
```

## Teardown

To delete everything:

```bash
# Delete Cloud Run service
gcloud run services delete recipebook-backend --region us-central1

# Destroy Terraform infrastructure
cd infra
terraform destroy

# Note: Database will be preserved due to force_destroy = false
# Manually delete bucket if needed:
# gsutil rm -r gs://$DB_BUCKET_NAME
```

## Security Considerations

### Production Checklist

- [ ] Remove `--allow-unauthenticated` and use Cloud Run IAM properly
- [ ] Set up CORS properly in backend
- [ ] Use Firebase App Check to prevent API abuse
- [ ] Enable Cloud Armor for DDoS protection
- [ ] Set up alerts for unusual activity
- [ ] Regular database backups (already handled by versioning)
- [ ] Review IAM permissions regularly

### Environment Variables

Never commit these to git:
- Service account JSON files
- Firebase config with API keys (use environment variables)
- Any passwords or secrets

## Next Steps

1. ✅ Backend deployed with `maxScale=1`
2. ✅ Database bucket created with versioning
3. ✅ Service account configured
4. ⏳ Build frontend with React
5. ⏳ Deploy frontend to Firebase Hosting
6. ⏳ Set up custom domain (optional)
7. ⏳ Add monitoring and alerts

## Support

See documentation:
- [DATABASE_ARCHITECTURE.md](DATABASE_ARCHITECTURE.md) - Database details
- [CLOUD_RUN_DATABASE_BEHAVIOR.md](CLOUD_RUN_DATABASE_BEHAVIOR.md) - How it works
- [backend/RECIPE_API.md](backend/RECIPE_API.md) - API documentation
