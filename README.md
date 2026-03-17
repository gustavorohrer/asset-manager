## Eclypsium BE Technical Challenge

This project implements a Go API for the Eclypsium backend challenge.

## TL;DR (Reviewer Path)
1. Run backend locally with seeded PostgreSQL (commands below).
2. Run the quick validation checklist (6 curl calls).
3. Run tests:
   - `go test ./...`
   - `DATABASE_URL="postgres://applicant:goodluck@localhost:5433/eclypsiumdb?sslmode=disable" go test -tags=integration ./integration`

## Demo
- Backend URL: `pending`
- Frontend URL: `pending`
- Video walkthrough: `pending`

Current implemented feature set:
- `GET /health`
- `GET /assets` (simple asset listing with filters, sorting, pagination, and computed threat/vulnerability flags)
- `GET /assets/:id` (asset details with ordered components and computed threat/vulnerability flags)
- `GET /assets/:id/vulnerabilities` (latest scan vulnerabilities by asset, with pagination and optional severity filter)
- `GET /assets/:id/threats` (latest scan threats by asset, with pagination and optional riskLevel filter)
- `PATCH /assets/:id` (partial update of asset fields: `name`, `description`, `lastScan`)
- `DELETE /assets/:id` (hard delete asset and related components/scans/vulnerabilities/threats)

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

## Reviewer Quick Validation Checklist
```bash
curl http://localhost:8080/health
curl "http://localhost:8080/assets?page=1&pageSize=3&sortBy=createdAt&sortOrder=desc"
curl "http://localhost:8080/assets/AST-001"
curl "http://localhost:8080/assets/AST-001/vulnerabilities?page=1&pageSize=5&severity=critical"
curl "http://localhost:8080/assets/AST-001/threats?page=1&pageSize=5&riskLevel=high"
curl -X PATCH "http://localhost:8080/assets/AST-001" -H "Content-Type: application/json" -d '{"name":"AST-001 Updated","description":"updated from reviewer checklist","lastScan":"2024-10-07T00:00:00Z"}'
```

Error contract checks:

```bash
curl "http://localhost:8080/assets?page=0"                        # 400 INVALID_QUERY_PARAM
curl "http://localhost:8080/assets/AST-404/threats"               # 404 ASSET_NOT_FOUND
curl -X PATCH "http://localhost:8080/assets/AST-001" -H "Content-Type: application/json" -d '{"id":"AST-002"}'  # 400 INVALID_REQUEST_BODY
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

### Asset vulnerabilities
```bash
curl "http://localhost:8080/assets/AST-001/vulnerabilities"
```

Supported query params:
- `page` (default `1`, max `10000`)
- `pageSize` (default `20`, max `100`)
- `severity` (`LOW`, `MEDIUM`, `HIGH`, `CRITICAL`, case-insensitive)

Example:

```bash
curl "http://localhost:8080/assets/AST-001/vulnerabilities?page=1&pageSize=10&severity=critical"
```

Success envelope:

```json
{
  "data": [
    {
      "id": "VUL-001",
      "description": "Dell UEFI BIOS vulnerable to BootHole (CVE-2020-10713) allowing SecureBoot bypass",
      "severity": "CRITICAL",
      "scanId": "SCN-001",
      "componentId": "CMP-001",
      "componentName": "Dell UEFI BIOS",
      "performedAt": "2024-10-08T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "total": 1,
    "totalPages": 1
  }
}
```

Invalid query example:

```bash
curl "http://localhost:8080/assets/AST-001/vulnerabilities?severity=urgent"
```

```json
{
  "error": {
    "code": "INVALID_QUERY_PARAM",
    "message": "one or more query parameters are invalid",
    "details": [
      {
        "field": "severity",
        "issue": "must be one of LOW, MEDIUM, HIGH, CRITICAL",
        "value": "urgent"
      }
    ]
  }
}
```

Not found example:

```bash
curl "http://localhost:8080/assets/AST-404/vulnerabilities"
```

```json
{
  "error": {
    "code": "ASSET_NOT_FOUND",
    "message": "asset not found"
  }
}
```

### Asset threats
```bash
curl "http://localhost:8080/assets/AST-001/threats"
```

Supported query params:
- `page` (default `1`, max `10000`)
- `pageSize` (default `20`, max `100`)
- `riskLevel` (`LOW`, `MEDIUM`, `HIGH`, case-insensitive)

Example:

```bash
curl "http://localhost:8080/assets/AST-001/threats?page=1&pageSize=10&riskLevel=high"
```

Success envelope:

```json
{
  "data": [
    {
      "id": "THR-001",
      "description": "Bootkits and rootkits can bypass SecureBoot via BootHole exploit",
      "riskLevel": "HIGH",
      "type": "Firmware Implant",
      "scanId": "SCN-001",
      "componentId": "CMP-001",
      "componentName": "Dell UEFI BIOS",
      "performedAt": "2024-10-08T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "total": 1,
    "totalPages": 1
  }
}
```

Invalid query example:

```bash
curl "http://localhost:8080/assets/AST-001/threats?riskLevel=critical"
```

```json
{
  "error": {
    "code": "INVALID_QUERY_PARAM",
    "message": "one or more query parameters are invalid",
    "details": [
      {
        "field": "riskLevel",
        "issue": "must be one of LOW, MEDIUM, HIGH",
        "value": "critical"
      }
    ]
  }
}
```

### Update asset (partial)
```bash
curl -X PATCH "http://localhost:8080/assets/AST-001" \
  -H "Content-Type: application/json" \
  -d '{"name":"AST-001 Updated","description":"updated description","lastScan":"2024-10-07T00:00:00Z"}'
```

Supported body fields:
- `name` (string, trimmed, non-empty, max 255)
- `description` (string, max 10000, empty string allowed)
- `lastScan` (RFC3339 string or `null` to clear value)

Rules:
- unknown body fields are rejected with `400 INVALID_REQUEST_BODY`
- `null` is not allowed for `name` and `description`
- at least one updatable field must be provided

Success envelope:

```json
{
  "data": {
    "id": "AST-001",
    "name": "AST-001 Updated",
    "description": "updated description",
    "createdAt": "2024-01-15T00:00:00Z",
    "lastScan": "2024-10-07T00:00:00Z"
  }
}
```

Invalid body example:

```bash
curl -X PATCH "http://localhost:8080/assets/AST-001" \
  -H "Content-Type: application/json" \
  -d '{"id":"AST-002"}'
```

```json
{
  "error": {
    "code": "INVALID_REQUEST_BODY",
    "message": "request body is invalid",
    "details": [
      {
        "field": "id",
        "issue": "is not allowed",
        "value": "\"AST-002\""
      }
    ]
  }
}
```

### Delete asset
Use this endpoint with care because it performs a hard delete.

```bash
curl -X DELETE "http://localhost:8080/assets/AST-001"
```

Success envelope:

```json
{
  "data": {
    "id": "AST-001",
    "deleted": true
  }
}
```

Not found example:

```bash
curl -X DELETE "http://localhost:8080/assets/AST-404"
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
- vulnerabilities query parsing/validation unit tests
- threats query parsing/validation unit tests
- update asset request parsing/validation unit tests
- service tests (`ListAssets`, `GetAssetDetails`, `ListAssetVulnerabilities`, `ListAssetThreats`, `UpdateAsset`, `DeleteAsset`)
- HTTP handler tests for `GET /assets`, `GET /assets/:id`, `GET /assets/:id/vulnerabilities`, `GET /assets/:id/threats`, `PATCH /assets/:id`, and `DELETE /assets/:id`
- integration tests against real PostgreSQL (`go test -tags=integration ./integration`)

Run integration tests:

```bash
DATABASE_URL="postgres://applicant:goodluck@localhost:5433/eclypsiumdb?sslmode=disable" go test -tags=integration ./integration
```

If `DATABASE_URL` is not set, integration tests are skipped.

## Evidence
- Unit and handler tests passing with `go test ./...`.
- Integration tests passing against real PostgreSQL with `-tags=integration`.
- Endpoints validated with success and error contracts (`200`, `400`, `404`) using seeded challenge data.

## Scope and Trade-offs
- Implemented core challenge endpoints: asset listing, details, vulnerabilities by asset, threats by asset, partial asset update, and hard delete.
- Kept API contract consistent (`{data, pagination}` for lists, structured error envelope).
- Used pgx + explicit SQL for clarity, control, and reproducibility in interview review.
- Prioritized atomic commits and test coverage (unit + integration) over advanced non-functional features.

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

## Backlog (Future Improvements)
- Auth/authz for write operations (`PATCH`, `DELETE`) and token-based access control.
- Optimistic concurrency control for `PATCH` (version/ETag + conditional update) to avoid lost updates.
- Audit trail for mutable operations (who/when/what changed, including deletes).
- OpenAPI/Swagger documentation and generated client examples.
- CI pipeline with automated unit/integration checks on pull requests.
