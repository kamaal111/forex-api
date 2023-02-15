start-db:
    gcloud emulators firestore start

run-dev:
    #!/bin/sh

    export GCP_PROJECT_ID=forex-api-daily
    export SERVER_ADDRESS=127.0.0.1:8000

    ~/go/bin/reflex -r '\.go' -s -- sh -c "go run ."

build:
    docker build -t forex-api .

run:
    #!/bin/sh

    export PORT=8000

    docker rm forex-api
    docker run -dp $PORT:$PORT --name forex-api -e GCP_PROJECT_ID=forex-api-daily -e PORT=$PORT forex-api
