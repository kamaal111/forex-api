import os
from google.cloud import firestore
from fastapi import FastAPI

app = FastAPI()

project_id = os.getenv("GCP_PROJECT_ID")
if not project_id:
    raise Exception("Failed to read GCP_PROJECT_ID environment variable")
db = firestore.AsyncClient(project=project_id)


@app.get("/exchange-rates")
async def list_exchange_rates():
    exchange_rates_collection = db.collection("exchange_rates")
    exchange_rates_documents = exchange_rates_collection.stream()
    exchange_rates = []
    async for exchange_rates_document in exchange_rates_documents:
        exchange_rates.append(exchange_rates_document.to_dict())
    return {"results": exchange_rates}
