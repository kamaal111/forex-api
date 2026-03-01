# Forex API

A lightweight REST API for retrieving foreign exchange rates, built with Go and backed by Google Cloud Firestore.

## Features

- 🚀 Fast and lightweight HTTP server using Go's standard library
- 💾 Firestore database for storing exchange rate data
- 🔄 Support for multiple base currencies
- 🎯 Filter rates by specific currency symbols
- 🐳 Docker support for easy deployment
- 📝 Request logging middleware

## Prerequisites

- Go 1.24+ (or Docker)
- Google Cloud Platform account with Firestore enabled
- [just](https://github.com/casey/just) command runner (optional, for development)
- [reflex](https://github.com/cespare/reflex) for hot-reloading during development
- Node.js 18+ and npm (for running integration tests via Firebase emulator)

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `GCP_PROJECT_ID` | Google Cloud Project ID with Firestore | Yes |
| `SERVER_ADDRESS` | Full server address (e.g., `127.0.0.1:8000`) | No |
| `PORT` | Port number (used if `SERVER_ADDRESS` not set) | Conditional |
| `FIRESTORE_EMULATOR_HOST` | Firestore emulator address for local development | No |

## Installation

### Clone the repository

```bash
git clone https://github.com/kamaal111/forex-api.git
cd forex-api
```

### Install dependencies

```bash
go mod download
```

## Usage

### Running Locally

#### Using just (recommended for development)

1. Start the Firestore emulator:
   ```bash
   just start-db
   ```

2. In a new terminal, start the development server with hot-reloading:
   ```bash
   just dev
   ```

#### Manual execution

```bash
export GCP_PROJECT_ID=your-project-id
export SERVER_ADDRESS=127.0.0.1:8000
go run .
```

### Running with Docker

#### Build the image

```bash
just build
# or
docker build -t forex-api .
```

#### Run the container

```bash
just run
# or
docker run -dp 8000:8000 --name forex-api \
  -e GCP_PROJECT_ID=your-project-id \
  -e PORT=8000 \
  forex-api
```

## API Endpoints

### Get All Available Currency Symbols

```
GET /v1/rates/symbols
```

Returns the full list of all supported currency symbols. Use this as a preflight check to discover which currencies are available before calling the rates endpoint.

#### Example Request

```bash
curl "http://localhost:8000/v1/rates/symbols"
```

#### Example Response

```json
["EUR","USD","JPY","BGN","CYP","CZK","DKK","EEK","GBP","HUF","LTL","LVL","MTL","PLN","ROL","RON","SEK","SIT","SKK","CHF","ISK","ILS","NOK","HRK","RUB","TRL","TRY","AUD","BRL","CAD","CNY","HKD","IDR","INR","KRW","MXN","MYR","NZD","PHP","SGD","THB","ZAR"]
```

### Get Latest Exchange Rates

```
GET /v1/rates/latest
```

Retrieves the most recent exchange rates for a given base currency.

#### Query Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| `base` | Base currency code (e.g., `USD`, `EUR`) | `EUR` |
| `symbols` | Comma-separated list of currency codes to filter (e.g., `USD,GBP,JPY`), or `*` to return all currencies | All currencies |

#### Example Request

```bash
curl "http://localhost:8000/v1/rates/latest?base=USD&symbols=EUR,GBP,JPY"
```

#### Example Response

```json
{
  "base": "USD",
  "date": "2025-11-30",
  "rates": {
    "EUR": 0.92,
    "GBP": 0.79,
    "JPY": 149.50
  }
}
```

### Supported Currencies

The API supports the following currencies:

- **Major**: EUR, USD, GBP, JPY, CHF, CAD, AUD, NZD
- **European**: BGN, CZK, DKK, HUF, PLN, RON, SEK, NOK, ISK, HRK
- **Asian**: CNY, HKD, IDR, INR, KRW, MYR, PHP, SGD, THB
- **Americas**: BRL, MXN
- **Other**: ILS, TRY, ZAR
- **Historical**: CYP, EEK, LTL, LVL, MTL, ROL, SIT, SKK, TRL

## Project Structure

```
forex-api/
├── main.go              # Application entry point
├── go.mod               # Go module dependencies
├── Dockerfile           # Docker container configuration
├── justfile             # Development task runner commands
├── database/
│   └── database.go      # Firestore client initialization
├── handlers/
│   └── rates.go         # HTTP request handlers for rates endpoint
├── routers/
│   ├── routers.go       # Main router setup and server start
│   ├── rates.go         # Rates route group
│   ├── middleware.go    # Request logging middleware
│   └── errors.go        # Error handling routes
└── utils/
    ├── environment.go   # Environment variable helpers
    ├── errors.go        # Error response utilities
    └── strings.go       # String utility functions
```

## Error Responses

All errors are returned in JSON format:

```json
{
  "message": "Error description",
  "status": 404
}
```

## Development

### Hot Reloading

The development setup uses [reflex](https://github.com/cespare/reflex) for automatic recompilation when Go files change:

```bash
go install github.com/cespare/reflex@latest
just dev
```

### Local Firestore Emulator

For local development without connecting to Google Cloud:

```bash
# Install Google Cloud SDK if not already installed
# https://cloud.google.com/sdk/docs/install

just start-db
```

This starts the Firestore emulator on `127.0.0.1:8080`.

## Testing

- Unit tests:
  - `just test` or `npm run test:unit`
- Integration tests (Firestore emulator managed automatically):
  - `just test-integration` or `npm test`
- All tests (unit + integration):
  - `just test-all` or `npm run test:all`
- Coverage:
  - `just test-cover` or `just test-cover-html`

Requirements for integration tests:
- Node.js 18+; `firebase-tools` is installed as a dev dependency and started by the npm scripts.
- No external GCP access is required; the tests seed data into the emulator.

## Contributing

See AGENTS.md for contributor guidelines, coding style, and PR expectations.
