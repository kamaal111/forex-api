import os
from typing import TYPE_CHECKING, List, Optional
from google.cloud import firestore
from fastapi import FastAPI, HTTPException

from constants import CURRENCIES


if TYPE_CHECKING:
    from internal_types import ExchangeRateRecord

app = FastAPI()

project_id = os.getenv("GCP_PROJECT_ID")
if not project_id:
    raise Exception("Failed to read GCP_PROJECT_ID environment variable")
db = firestore.AsyncClient(project=project_id)


@app.get("/v1/latest")
async def get_latest(base: Optional[str] = None, symbols: Optional[str] = None):
    if base is None or base not in CURRENCIES:
        base = "EUR"
    else:
        base = base.upper()

    exchange_rates_result = (
        db.collection("exchange_rates")
        .order_by("date", direction=firestore.Query.DESCENDING)
        .where("base", "==", base)
        .limit(1)
        .stream()
    )
    exchange_rate: Optional[ExchangeRateRecord] = None
    async for item in exchange_rates_result:
        exchange_rate = item.to_dict()
        break  # shouldn't loop more than once anyway, but let's make sure we can sleep at night

    if exchange_rate is None:
        raise HTTPException(status_code=404, detail="Item not found")

    symbols_list = make_symbols_list(raw=symbols, base=base)
    if len(symbols_list) == 0:
        return exchange_rate

    exchange_rate_with_filtered_rates: ExchangeRateRecord = {
        **exchange_rate,
        "rates": {},
    }
    for symbol in symbols_list:
        exchange_rate_with_filtered_rates["rates"][symbol] = exchange_rate["rates"][
            symbol
        ]
    return exchange_rate_with_filtered_rates


def make_symbols_list(raw: Optional[str], base: str):
    if not raw:
        return []

    symbols_list: List[str] = []
    for item in raw.upper().split(","):
        if item != base and item in CURRENCIES:
            symbols_list.append(item)

    return symbols_list
