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
