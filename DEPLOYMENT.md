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
export PROJECT_ID="projectid"

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

### 6. Configure Workload Identity Federation

#### Create Workload Identity Pool

```bash
gcloud iam workload-identity-pools create "github-actions-pool" \
    --project="$PROJECT_ID" \
    --location="global" \
    --display-name="GitHub Actions Pool"
```

#### Create Workload Identity Provider

```bash
gcloud iam workload-identity-pools providers create-oidc "github-actions-provider" \
    --project="$PROJECT_ID" \
    --location="global" \
    --workload-identity-pool="github-actions-pool" \
    --display-name="GitHub Actions Provider" \
    --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
    --attribute-condition="assertion.repository_owner=='nischith77'" \
    --issuer-uri="https://token.actions.githubusercontent.com"
```

#### Bind Service Account to Workload Identity

```bash
# Replace YOUR_GITHUB_REPO with your GitHub repository (e.g., "nischith77/go-vehicle-detail")
export GITHUB_REPO="repo"

gcloud iam service-accounts add-iam-policy-binding \
    --project="$PROJECT_ID" \
    --role="roles/iam.workloadIdentityUser" \
    --member="principalSet://iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/github-actions-pool/attribute.repository/$GITHUB_REPO" \
    "github-actions@$PROJECT_ID.iam.gserviceaccount.com"
```

#### Get WIF Provider Resource Name

```bash
gcloud iam workload-identity-pools providers describe "github-actions-provider" \
    --project="$PROJECT_ID" \
    --location="global" \
    --workload-identity-pool="github-actions-pool" \
    --format="value(name)"
```

Save the output (WIF Provider resource name) for GitHub secrets configuration.

## GitHub Secrets Configuration

Add the following secrets to your GitHub repository (Settings → Secrets and variables → Actions):

### Required Secrets

1. **WIF_PROVIDER**: Workload Identity Provider resource name from step 6 above
   - Format: `projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/github-actions-pool/providers/github-actions-provider`
2. **WIF_SERVICE_ACCOUNT**: Service account email
   - Format: `github-actions@YOUR_PROJECT_ID.iam.gserviceaccount.com`
3. **GCP_PROJECT_ID**: Your GCP project ID
4. **CLOUD_SQL_INSTANCE**: Cloud SQL instance connection name
   - Format: `PROJECT_ID:REGION:INSTANCE_NAME` (e.g., `mobility-tracker-467907:us-central1:mobility-tracker`)
5. **DB_CONN**: PostgreSQL connection string for your Cloud SQL instance
   - For Unix socket (recommended for Cloud Run): `host=/cloudsql/mobility-tracker-467907:us-central1:mobility-tracker user=postgres password=your-password dbname=vehicledb sslmode=disable`
   - For TCP (alternative): `host=PRIVATE_IP port=5432 user=postgres password=your-password dbname=vehicledb sslmode=require`

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

- **Enhanced Security**: Uses Workload Identity Federation (WIF) instead of long-lived service account keys
- **Automatic Token Rotation**: Authentication tokens are automatically managed and rotated by Google Cloud
- **Reduced Attack Surface**: No sensitive service account keys stored in GitHub secrets
- **Principle of Least Privilege**: WIF allows fine-grained access control based on repository and branch
- Use Cloud SQL Proxy or Private IP for database connections in production
- Consider restricting WIF bindings to specific branches for additional security
