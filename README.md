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


## Authentication (JWT & Refresh Tokens)
This backend uses a dual-token system for security:
- **Access Token**: Short-lived JWT (15-60 mins) for API authorization.
- **Refresh Token**: Long-lived opaque token (10h+) stored in the database.
- **Token Rotation**: Every time a session is refreshed, a new refresh token is issued and the old one is revoked.
- **Theft Detection**: Reusing a revoked refresh token will trigger an account-wide session revocation for that user.

## Background Jobs (Cron)
The server runs scheduled tasks using `robfig/cron`:
- **Daily Cleanup**: Every day at 03:00 Bangkok time, it deletes expired and revoked refresh tokens from the database.
- **Reminders**: Sends periodic email notifications for pending tasks.

## Running Backend Locally
1. Place `storageLocal.json` and create `.env` file from `.env.example`.
2. Run the server: `go run cmd/api-server/main.go`
3. The server will start on `:8080` (default).
4. Log in via `POST /auth/login` to get your initial tokens.

## Development & Quality Control
To ensure code quality and prevent CI failures, please run linting and tests locally before creating a Pull Request.

### Prerequisites
1. **Go 1.24.3**
2. **golangci-lint**: 
   - Windows: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
   - Others: [Installation Guide](https://golangci-lint.run/usage/install/)
3. **C compiler (Optional)**: Needed for the `-race` flag. If you don't have one, our script will automatically run tests without it.

### Quick Check Script
We provide scripts to run all checks in one command:

**Windows (PowerShell):**
```powershell
.\scripts\check.ps1
```

**macOS/Linux (Bash):**
```bash
./scripts/check.sh
```

### Using Makefile (macOS/Linux)
- **Run all checks**: `make all`
- **Run Unit Tests**: `make test`
- **Run Linter**: `make lint`

### Manual Commands (Cross-platform)
If you don't have `make` or want to run specific steps:
- **Run Everything**: `golangci-lint run; go test -v -race ./...`
- **Unit Tests Only**: `go test -p 4 -v -race ./...`
- **Linter Only**: `golangci-lint run`

### PR Workflow
1. Develop features and add unit tests.
2. Run the check script for your platform:
   - Windows: `.\scripts\check.ps1`
   - macOS/Linux: `./scripts/check.sh` or `make all`
3. Push changes and create a PR (CI will verify again).

<!-- Deploy on cloud -->
gcloud run deploy trl-research-backend \
  --source . \
  --region asia-southeast1 \
  --allow-unauthenticated

gcloud run deploy trl-research-backend --source . --region asia-southeast1 --allow-unauthenticated