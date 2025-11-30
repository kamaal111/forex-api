start-db:
    gcloud emulators firestore start

dev:
    #!/bin/sh

    export GCP_PROJECT_ID=forex-api-daily
    export SERVER_ADDRESS=127.0.0.1:8000
    export FIRESTORE_EMULATOR_HOST="127.0.0.1:8080"

    ~/go/bin/reflex -r '\.go' -s -- sh -c "go run ."

build:
    docker build -t forex-api .

run:
    #!/bin/sh

    export PORT=8000

    docker stop forex-api || true
    docker rm forex-api || true
    docker run -dp $PORT:$PORT --name forex-api -e FIRESTORE_EMULATOR_HOST=host.docker.internal:8080 \
        -e GCP_PROJECT_ID=forex-api-daily -e SERVER_ADDRESS=0.0.0.0:$PORT forex-api

build-run: build run
