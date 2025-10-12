// build
gcloud builds submit --tag asia-southeast1-docker.pkg.dev/gen-lang-client-0264595198/trl-research-backend-repo/trl-research-backend

// deploy
gcloud run deploy trl-research-backend \
   --image asia-southeast1-docker.pkg.dev/gen-lang-client-0264595198/trl-research-backend-repo/trl-research-backend \
   --platform managed \
   --region asia-southeast1 \
   --allow-unauthenticated

// url
Service URL: https://trl-research-backend-325350196988.asia-southeast1.run.app


## run backend on local
1. put file trl-research-service-account.json, private_key_v1.pem, public_key_v1.pem on same layer as Dckerfile

2. go run cmd/api-server/main.go

3. login before call any API - looking for admin account in file internal/script/seed_admins.go

<!-- Deploy on cloud -->
gcloud run deploy trl-research-backend \
  --source . \
  --region asia-southeast1 \
  --allow-unauthenticated
