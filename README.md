# GCP SERVERLESS API

## TECH STACK

- API Gateway
- Golang
- Google Cloud Firestore (Native Mode)
- Google Cloud Functions
- Google Cloud Storage

## LOCAL SETUP

- `git clone <THIS REPO URL>` and `cd` into the folder
- `go mod download` to download modules to local cache
- Create a service account and private key for the account, add roles `Cloud Functions Invoker`,`Storage Object Admin`, `Cloud Datastore User`
- Create a publicly accessible Google Cloud Storage Bucket
- Copy the following environmental variables to a `.env` file from `example.env` and fill in your credentials
- `make dev` OR `go run cmd/main.go` to start the application

## GCP DEPLOYMENT

- Install the [serverless framework](https://www.serverless.com/framework/docs/getting-started) and plugins defined in `serverless.yml`
- Enable these APIs on GCP
  - `Cloud Deployment Manager V2 API`,
  - `API Gateway API`,
  - `Service Control API`,
  - `Service Management API`,
  - `Cloud Firestore API (Native Mode)`,
  - `Cloud Function API`
  - `Cloud Storage API`
- Use gcloud CLI to get application default credientials
- For each stage e.g. dev, prod, test
  - Copy the following environmental variables to a `.env.<STAGE_NAME>` file from `example.env` and fill in your credentials
  - `sls deploy --stage <STAGE_NAME>` to deploy resources
  - Create [an API Gateway and its API config](https://cloud.google.com/api-gateway/docs/creating-api?hl=en) using the sample file `apiconfig-dev.yaml`. Replace the region and project id in the given file
  - Enable your created API e.g. <API_ID>-<HASH>.apigateway.<PROJECT_ID>.cloud.goog
  - Create an API Key and restrict it to your API
