## Eclypsium BE Technical Challenge

This project implements a Go API for the Eclypsium backend challenge.

Current implemented feature set:
- `GET /health`
- `GET /assets` (simple asset listing with filters, sorting, pagination, and computed threat/vulnerability flags)
- `GET /assets/:id` (asset details with ordered components and computed threat/vulnerability flags)

## Tech stack
- Go 1.25
- Gin (HTTP server)
- PostgreSQL 17
- pgx v5 (SQL access)

## Reviewer quick start (copy/paste)

This path avoids conflicts with a local PostgreSQL by using host port `5433`.

```bash
docker build -t ecl-be-challenge-db ./db
docker rm -f ecl-be-challenge-db 2>/dev/null || true
docker run --name ecl-be-challenge-db \
  -e POSTGRES_DB=eclypsiumdb \
  -e POSTGRES_USER=applicant \
  -e POSTGRES_PASSWORD=goodluck \
  -p 5433:5432 \
  -d ecl-be-challenge-db

DATABASE_URL="postgres://applicant:goodluck@localhost:5433/eclypsiumdb?sslmode=disable" go run .
```

In another terminal:

```bash
curl http://localhost:8080/health
curl "http://localhost:8080/assets?page=1&pageSize=5&sortBy=createdAt&sortOrder=desc"
```

If you prefer separate env vars instead of `DATABASE_URL`:
- `DB_HOST=localhost`
- `DB_PORT=5432` (or `5433` if using the quick start above)
- `DB_NAME=eclypsiumdb`
- `DB_USER=applicant`
- `DB_PASSWORD=goodluck`
- `PORT=8080` (optional, default is `8080`)

You can also set `DATABASE_URL` directly (it takes precedence).

## API

### Health check
```bash
curl http://localhost:8080/health
```

### Asset listing
```bash
curl "http://localhost:8080/assets"
```

Supported query params:
- `name` (contains, case-insensitive)
- `created_from` (RFC3339)
- `created_to` (RFC3339)
- `last_scan_from` (RFC3339)
- `last_scan_to` (RFC3339)
- `page` (default `1`, max `10000`)
- `pageSize` (default `20`, max `100`)
- `sortBy` (`createdAt`, `name`, `lastScan`)
- `sortOrder` (`asc`, `desc`)

Examples:

```bash
curl "http://localhost:8080/assets?page=1&pageSize=10&sortBy=name&sortOrder=asc"
```

```bash
curl "http://localhost:8080/assets?name=router&created_from=2024-01-01T00:00:00Z&created_to=2024-12-31T23:59:59Z"
```

Invalid query example:

```bash
curl "http://localhost:8080/assets?page=0"
```

Returns:

```json
{
  "error": {
    "code": "INVALID_QUERY_PARAM",
    "message": "one or more query parameters are invalid",
    "details": [
      {
        "field": "page",
        "issue": "must be a positive integer",
        "value": "0"
      }
    ]
  }
}
```

Listing success envelope:

```json
{
  "data": [
    {
      "id": "AST-001",
      "name": "Dell PowerEdge R740 Server",
      "description": "Production database server in datacenter rack A3",
      "createdAt": "2024-01-15T00:00:00Z",
      "lastScan": "2024-10-08T00:00:00Z",
      "hasVulnerabilities": true,
      "hasThreats": true
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 20,
    "total": 12,
    "totalPages": 1
  }
}
```

### Asset details
```bash
curl "http://localhost:8080/assets/AST-001"
```

Success envelope:

```json
{
  "data": {
    "id": "AST-001",
    "name": "Dell PowerEdge R740 Server",
    "description": "Production database server in datacenter rack A3",
    "createdAt": "2024-01-15T00:00:00Z",
    "lastScan": "2024-10-08T00:00:00Z",
    "hasVulnerabilities": true,
    "hasThreats": true,
    "components": [
      {
        "id": "CMP-001",
        "name": "Dell UEFI BIOS",
        "version": "2.10.2",
        "vendor": "Dell Inc.",
        "type": "UEFI",
        "createdAt": "2024-01-15T00:00:00Z",
        "lastScan": "2024-10-08T00:00:00Z",
        "assetId": "AST-001"
      }
    ]
  }
}
```

Not found example:

```bash
curl "http://localhost:8080/assets/AST-404"
```

```json
{
  "error": {
    "code": "ASSET_NOT_FOUND",
    "message": "asset not found"
  }
}
```

## Tests
Run:

```bash
go test ./...
```

Current tests include:
- query parsing/validation unit tests
- service tests (`ListAssets`, `GetAssetDetails`)
- HTTP handler tests for `GET /assets` and `GET /assets/:id`

## Notes
- Unknown query params are ignored.
- Date filters are applied against database `DATE` columns using the date portion of the RFC3339 value.
- `lastScan` null assets are excluded only when `last_scan_from` or `last_scan_to` filters are present.
- For `sortBy=lastScan`, `NULLS LAST` is applied.

## Troubleshooting
- Error `role "applicant" does not exist`:
  - You are probably connecting to a different local PostgreSQL instance.
  - Use the quick start above (`5433` + `DATABASE_URL`) to force the API to use the challenge DB container.
- Port already in use:
  - Change host port in `docker run` (for example `-p 5440:5432`) and update `DATABASE_URL`.

## Cleanup
```bash
docker rm -f ecl-be-challenge-db
```

## Backlog
- `GET /assets/:id/vulnerabilities` (latest scan logic per component)
- `GET /assets/:id/threats` (latest scan logic per component)
- Optional challenge endpoints:
  - update asset properties
  - remove asset
- Integration tests against real PostgreSQL
