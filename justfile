start-db:
    gcloud emulators firestore start

run:
    #!/bin/sh

    export GCP_PROJECT_ID=forex-api

    uvicorn main:app --reload
