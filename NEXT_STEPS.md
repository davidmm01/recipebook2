# Next Steps for Deployment

## 1. Run the Deploy Workflow

Go to: GitHub repo > Actions > Deploy > Run workflow

You can choose to deploy:
- Backend only
- Frontend only
- Both (default)

## 2. Verify

- Backend: https://recipebook-backend-764030971451.australia-southeast2.run.app/health
- Frontend: Check Firebase Hosting URL (will be shown in workflow output)
- Backups: Check Cloud Storage bucket after 6 hours for first backup
