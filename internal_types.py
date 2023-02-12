from typing import Dict, TypedDict


class ExchangeRateRecord(TypedDict):
    base: str
    date: str
    rates: Dict[str, float]
