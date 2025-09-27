This is standalone go app which reads vehicle details from https://carapi.app/api/models/v2 as json and saves it to table in postgres.

## Local Development

### Steps to run locally
1. install docker desktop
2. cd to go-vehicle-detail path
3. run docker-compose up

## Cloud Deployment

This application can be deployed to Google Cloud Platform as a Cloud Run Job using GitHub Actions.

### Automated Deployment
- Push to `master` or `main` branch triggers automatic deployment
- Manual deployment available via GitHub Actions workflow dispatch

### Setup Instructions
See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed setup instructions including:
- GCP project configuration
- Service account setup
- GitHub secrets configuration
- Database setup

### Running on GCP
Once deployed, the job can be executed manually or scheduled:
```bash
gcloud run jobs execute go-vehicle-detail-job --region=us-central1
```
