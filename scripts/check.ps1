# Check for golangci-lint
if (!(Get-Command golangci-lint -ErrorAction SilentlyContinue)) {
    Write-Host "WARNING: golangci-lint is not installed." -ForegroundColor Yellow
    Write-Host "To install it, run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" -ForegroundColor Cyan
    Write-Host "Skipping linting step..." -ForegroundColor Gray
    Write-Host ""
}
else {
    Write-Host "Running Linter..." -ForegroundColor Cyan
    golangci-lint run
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Linting failed!" -ForegroundColor Red
        exit $LASTEXITCODE
    }
}

Write-Host "Running Unit Tests..." -ForegroundColor Cyan
# Skip -race if CGO is not enabled (common on Windows without GCC)
$goEnv = go env CGO_ENABLED
if ($goEnv -eq "1") {
    go test -v -race ./...
}
else {
    Write-Host "INFO: CGO is disabled, running tests without -race flag" -ForegroundColor Gray
    go test -v ./...
}

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Tests failed!" -ForegroundColor Red
    exit $LASTEXITCODE
}

Write-Host "SUCCESS: All checks passed!" -ForegroundColor Green
