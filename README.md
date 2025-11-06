# RecipeBook

A simple web application for managing recipes for food, cocktails, and more.

## Stack

**Frontend:**
- React (hosted on Firebase Hosting)

**Backend:**
- Go REST API (deployed on Google Cloud Run)

**Database:**
- SQLite + Cloud Storage (see [docs/DATABASE_ARCHITECTURE.md](docs/DATABASE_ARCHITECTURE.md) for details)

**Infrastructure:**
- Google Cloud Platform (GCP)
- Terraform for infrastructure as code

## Project Structure

```
recipebook2/
├── frontend/       # React application
├── backend/        # Go REST API
├── infra/          # Terraform infrastructure code
└── README.md
```

## Features

- **Public recipe viewing** - Anyone can browse and search recipes
- **Authenticated editing** - Only you (and your girlfriend) can create/edit/delete
- Markdown formatting for ingredients and methods
- Full-text search across all recipes
- Support for multiple recipe types (food, cocktails, etc.)
- Firebase Authentication for editors (email/password)
- Ultra-low cost (~$0.001/month)

## Documentation

📚 **[See docs/ for complete documentation](docs/)** including:
- Deployment guide
- Database architecture
- Security model
- API reference

## Quick Start

For complete deployment instructions, see **[docs/DEPLOYMENT.md](docs/DEPLOYMENT.md)**.

## Local Development Setup

### 1. Firebase Project Setup

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Create a new project (or select existing)
3. Enable **Authentication**:
   - Go to Authentication > Sign-in method
   - Enable "Email/Password"
4. Get your Firebase config:
   - Go to Project Settings > General
   - Scroll to "Your apps" and click web icon to register app
   - Copy the Firebase configuration

### 2. Frontend Setup

```bash
cd frontend

# Copy environment template
cp .env.example .env

# Edit .env and add your Firebase config values
# Get these from Firebase Console > Project Settings > General > Your apps

# Install dependencies (already done if you followed initial setup)
npm install

# Run development server
npm start
```

### 3. Infrastructure Setup

```bash
cd infra

# Create terraform.tfvars file
echo 'project_id = "your-gcp-project-id"' > terraform.tfvars

# Initialize and apply
terraform init
terraform plan
terraform apply

# Note the output values (bucket name and service account)
```

### 4. Backend Setup

```bash
cd backend

# For local development, download service account key:
# Firebase Console > Project Settings > Service Accounts > Generate new private key
# Save as service-account.json in backend/ directory

# Set environment variables
export GOOGLE_APPLICATION_CREDENTIALS="./service-account.json"
export DB_BUCKET_NAME="your-project-id-recipebook-db"

# Run server
go run .
```

See [docs/DATABASE_ARCHITECTURE.md](docs/DATABASE_ARCHITECTURE.md) for detailed database setup and deployment instructions.

## Development

**Frontend:** http://localhost:3000
**Backend API:** http://localhost:8080

## Authentication Flow

1. User signs up/logs in via frontend (Firebase Auth)
2. Frontend receives Firebase ID token
3. Frontend sends token in `Authorization: Bearer <token>` header to backend
4. Backend validates token with Firebase Admin SDK
5. If valid, request proceeds; otherwise returns 401

For detailed API documentation, see **[docs/API.md](docs/API.md)**.
