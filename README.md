# RecipeBook

A simple web application for managing recipes for food, cocktails, and more.

## Stack

**Frontend:**
- React (hosted on Firebase Hosting)

**Backend:**
- Go REST API (deployed on Google Cloud Run)

**Database:**
- Firestore (NoSQL document store)

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

- View recipes
- Create new recipes
- Update existing recipes
- Support for multiple recipe types (food, cocktails, etc.)
- Firebase Authentication (email/password)

## Setup

### 1. Firebase Project Setup

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Create a new project (or select existing)
3. Enable **Authentication**:
   - Go to Authentication > Sign-in method
   - Enable "Email/Password"
4. Enable **Firestore Database**:
   - Go to Firestore Database > Create database
   - Start in production mode (we'll add security rules later)
5. Get your Firebase config:
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

### 3. Backend Setup

```bash
cd backend

# For local development, download service account key:
# Firebase Console > Project Settings > Service Accounts > Generate new private key
# Save as service-account.json in backend/ directory

# Set environment variable
export GOOGLE_APPLICATION_CREDENTIALS="./service-account.json"

# Run server
go run .
```

### 4. Infrastructure Setup

```bash
cd infra

# Create terraform.tfvars file
echo 'project_id = "your-gcp-project-id"' > terraform.tfvars

# Initialize and apply
terraform init
terraform plan
terraform apply
```

## Development

**Frontend:** http://localhost:3000
**Backend API:** http://localhost:8080

## Authentication Flow

1. User signs up/logs in via frontend (Firebase Auth)
2. Frontend receives Firebase ID token
3. Frontend sends token in `Authorization: Bearer <token>` header to backend
4. Backend validates token with Firebase Admin SDK
5. If valid, request proceeds; otherwise returns 401

## API Endpoints

- `GET /health` - Health check (no auth required)
- `GET /recipes` - List recipes (auth required)
- `POST /recipes` - Create recipe (auth required)
