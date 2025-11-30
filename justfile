# Start the Firestore emulator manually
start-db:
    gcloud emulators firestore start

# Run the development server with hot reload
dev:
    #!/bin/sh

    export GCP_PROJECT_ID=forex-api-daily
    export SERVER_ADDRESS=127.0.0.1:8000
    export FIRESTORE_EMULATOR_HOST="127.0.0.1:8080"

    ~/go/bin/reflex -r '\.go' -s -- sh -c "go run ."

# Run unit tests only (no emulator needed)
test:
    go test ./... -v -short

# Run unit tests with coverage report
test-cover:
    go test ./... -cover -short

# Run unit tests with coverage and open HTML report
test-cover-html:
    go test ./... -coverprofile=coverage.out -short && go tool cover -html=coverage.out

# Run integration tests with Firestore emulator (auto-starts and cleans up)
test-integration:
    pnpm test

# Run all tests including integration tests with Firestore emulator
test-all:
    pnpm run test:all

# Build the Docker image
build:
    docker build -t forex-api .

# Run the Docker container
run:
    #!/bin/sh

    export PORT=8000

    docker stop forex-api || true
    docker rm forex-api || true
    docker run -dp $PORT:$PORT --name forex-api -e FIRESTORE_EMULATOR_HOST=host.docker.internal:8080 \
        -e GCP_PROJECT_ID=forex-api-daily -e SERVER_ADDRESS=0.0.0.0:$PORT forex-api

# Build and run the Docker container
build-run: build run
