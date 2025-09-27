# GCP Cloud Run Job Deployment Setup

This document explains how to set up GCP Cloud Run Job deployment for the Go Vehicle Detail application using GitHub Actions.

## Prerequisites

1. Google Cloud Platform (GCP) account
2. A GCP project with billing enabled
3. GitHub repository with admin access

## GCP Setup

### 1. Enable Required APIs

```bash
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable artifactregistry.googleapis.com
gcloud services enable sql-component.googleapis.com
```

### 2. Create Artifact Registry Repository

```bash
gcloud artifacts repositories create go-vehicle-detail \
    --repository-format=docker \
    --location=us-central1 \
    --description="Repository for go-vehicle-detail app"
```

### 3. Create Cloud SQL PostgreSQL Instance (if needed)

```bash
gcloud sql instances create vehicle-detail-db \
    --database-version=POSTGRES_15 \
    --tier=db-f1-micro \
    --region=us-central1

gcloud sql databases create vehicledb --instance=vehicle-detail-db

gcloud sql users create postgres \
    --instance=vehicle-detail-db \
    --password=Nis@5020
```

### 4. Create Service Account

```bash
gcloud iam service-accounts create github-actions \
    --description="Service account for GitHub Actions" \
    --display-name="GitHub Actions"
```

### 5. Grant Required Permissions

```bash
# Replace YOUR_PROJECT_ID with your actual project ID
export PROJECT_ID="YOUR_PROJECT_ID"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/run.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/storage.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:github-actions@$PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser"
```

### 6. Create and Download Service Account Key

```bash
gcloud iam service-accounts keys create key.json \
    --iam-account=github-actions@$PROJECT_ID.iam.gserviceaccount.com
```

## GitHub Secrets Configuration

Add the following secrets to your GitHub repository (Settings → Secrets and variables → Actions):

### Required Secrets

1. **GCP_SA_KEY**: Content of the `key.json` file created above
2. **GCP_PROJECT_ID**: Your GCP project ID
3. **DB_CONN**: PostgreSQL connection string for your Cloud SQL instance
   - Format: `host=INSTANCE_IP port=5432 user=appuser password=your-secure-password dbname=vehicledb sslmode=require`

### Optional Secrets (with defaults)

1. **GAR_LOCATION**: Artifact Registry location (default: `us-central1`)
2. **GAR_REPOSITORY**: Artifact Registry repository name (default: `go-vehicle-detail`)
3. **CLOUD_RUN_REGION**: Cloud Run region (default: `us-central1`)

## Database Setup

If using Cloud SQL, you'll need to:

1. Create the database schema by running the `init.sql` file:

   ```bash
   gcloud sql connect vehicle-detail-db --user=appuser
   # Then run the contents of init.sql
   ```

2. Configure the connection string in the `DB_CONN` secret.

## Deployment

The GitHub Actions workflow will automatically trigger on:

- Push to `master` or `main` branch
- Manual workflow dispatch

The workflow will:

1. Build the Docker image
2. Push to Google Artifact Registry
3. Deploy as a Cloud Run Job
4. Set the required environment variables

## Running the Job

After deployment, you can execute the job manually:

```bash
gcloud run jobs execute go-vehicle-detail-job --region=us-central1
```

## Monitoring

You can monitor job executions in the GCP Console:

- Cloud Run → Jobs → go-vehicle-detail-job
- View logs and execution history

## Security Notes

- The service account key should be kept secure and rotated regularly
- Use Cloud SQL Proxy or Private IP for database connections in production
- Consider using Workload Identity Federation instead of service account keys for enhanced security
