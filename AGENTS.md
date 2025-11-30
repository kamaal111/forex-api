# Repository Guidelines

## Project Structure & Module Organization
- `main.go`: entry point.
- `routers/`: HTTP routing, middleware, errors.
- `handlers/`: request handlers and domain service (`service.go`).
- `database/`: Firestore client creation.
- `utils/`: env, errors, helpers.
- `test/integration/`: end-to-end tests using Firestore emulator.

## Build, Test, and Development Commands
- Local dev with hot reload:
  - `just start-db` — start Firestore emulator.
  - `just dev` — run server with reflex.
- Manual run:
  - `GCP_PROJECT_ID=... SERVER_ADDRESS=127.0.0.1:8000 go run .`
- Tests:
  - `just test` or `npm run test:unit` — unit tests only.
  - `just test-integration` or `npm test` — integration tests via emulator.
  - `just test-all` or `npm run test:all` — unit + integration.
- Docker:
  - `just build` then `just run` (maps port and sets envs).

## Coding Style & Naming Conventions
- Language: Go 1.24+ (module targets `go 1.25`).
- Formatting: `gofmt` defaults (tabs, goimports-style grouping). Run `go fmt ./...`.
- Naming: exported `CamelCase`, unexported `lowerCamel`; package names short lowercase.
- Files: use `snake_case.go`; tests end with `_test.go`.

## Testing Guidelines
- Framework: standard `testing` package.
- Unit tests live next to code (e.g., `handlers/service_test.go`). Name tests `TestXxx`.
- Integration tests in `test/integration/`; they start the app and hit `/v1/rates/latest`.
- Emulator envs: set `FIRESTORE_EMULATOR_HOST=127.0.0.1:8080` and `GCP_PROJECT_ID`.
- Prefer deterministic data: seed `exchange_rates` docs within tests.

## Commit & Pull Request Guidelines
- Commit style (observed): short, imperative summaries without prefixes.
  - Examples: "Connect to database", "Upgrade dependencies".
- PRs should include:
  - Summary of changes and rationale.
  - Test coverage notes (unit/integration) and run commands.
  - Any API changes with example `curl` requests/responses.
  - Linked issue(s) when applicable.

## Security & Configuration Tips
- Do not commit secrets; use env vars: `GCP_PROJECT_ID`, `SERVER_ADDRESS`, `PORT`, `FIRESTORE_EMULATOR_HOST`.
- Use the emulator for local/integration runs; real GCP only in trusted environments.

## Architecture Overview
- Minimal net/http server; routes in `routers`, logic in `handlers/service.go`, Firestore via `database`.
- Primary endpoint: `GET /v1/rates/latest?base=EUR&symbols=USD,GBP`.
