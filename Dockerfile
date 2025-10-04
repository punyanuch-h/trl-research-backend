# --------------------------
# STEP 1: Build Go binary
# --------------------------
    FROM golang:1.24 AS builder

    # Set working directory
    WORKDIR /app
    
    # Copy go.mod and go.sum
    COPY go.mod go.sum ./

    # ✅ เปิดใช้ auto toolchain
    ENV GOTOOLCHAIN=auto
    
    # Download dependencies
    RUN go mod download
    
    # Copy source code
    COPY . .
    
    # Build Go binary
    RUN go build -o server ./cmd/api-server/main.go
    
    # --------------------------
    # STEP 2: Run binary
    # --------------------------
    FROM gcr.io/distroless/base-debian12
    
    WORKDIR /app
    
    # Copy binary from builder
    COPY --from=builder /app/server .
    
    # Copy Firebase service account key
    COPY serviceAccountKey.json .
    
    # Expose port
    EXPOSE 8080
    
    # Run server
    CMD ["./server"]
    