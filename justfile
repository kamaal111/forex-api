start-db:
    gcloud emulators firestore start

run-dev:
    #!/bin/sh

    export GCP_PROJECT_ID=forex-api

    uvicorn main:app --reload

build:
    docker build -t forex-api .

run:
    docker run -dp 8000:8000  --name forex-api forex-api
