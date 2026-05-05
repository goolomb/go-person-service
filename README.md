# Go Person Service

Small REST API service for saving and reading people by external ID.

## Tech Stack

- Go
- chi
- GORM
- SQLite

## API

| Method | Path | Description | Expected status codes |
|---|---|---|---|
| GET | `/health` | Health check. | `200` |
| POST | `/save` | Validate and save a person. | `201`, `400`, `409`, `500` |
| GET | `/{id}` | Read a person by `external_id`. | `200`, `400`, `404`, `500` |

`POST /save` accepts:

```json
{
  "external_id": "3f93df6d-ff51-4740-9d27-fc6b2f30281c",
  "name": "Jane Doe",
  "email": "jane@example.com",
  "date_of_birth": "1990-01-02T03:04:05Z"
}
```

Successful `POST /save` requests return `201 Created`, a response body, and a `Location` header set to `/{external_id}`.

Successful `POST /save` and `GET /{id}` responses return:

```json
{
  "external_id": "3f93df6d-ff51-4740-9d27-fc6b2f30281c",
  "name": "Jane Doe",
  "email": "jane@example.com",
  "date_of_birth": "1990-01-02T03:04:05Z"
}
```

`GET /{id}` uses the person's `external_id`.

## Error Responses

General errors return:

```json
{
  "error": "not_found",
  "message": "Person was not found."
}
```

Validation errors return:

```json
{
  "error": "validation_error",
  "message": "Request body contains validation errors.",
  "fields": {
    "email": "email must be a valid email address"
  }
}
```

## Configuration

| Variable | Default | Description |
|---|---|---|
| `HTTP_PORT` | `8080` | Port where the HTTP server listens. |
| `DB_PATH` | `./data/app.db` | Path to the SQLite database file. |

## Run Locally

```sh
go run ./cmd/api
```

With custom configuration:

```sh
HTTP_PORT=9090 DB_PATH=/tmp/person-service.db go run ./cmd/api
```

## Build Binary

```sh
go build -o person-service ./cmd/api
./person-service
```

## Docker

Build:

```sh
docker build -t go-person-service .
```

Run:

```sh
docker run --rm -p 8080:8080 go-person-service
```

Run with persistent SQLite data:

```sh
docker run --rm \
  -p 8080:8080 \
  -v "$PWD/data:/app/data" \
  go-person-service
```

The container defaults to `HTTP_PORT=8080` and `DB_PATH=/app/data/app.db`.

## curl Examples

Save a person:

```sh
curl -i -X POST http://localhost:8080/save \
  -H 'Content-Type: application/json' \
  -d '{
    "external_id": "3f93df6d-ff51-4740-9d27-fc6b2f30281c",
    "name": "Jane Doe",
    "email": "jane@example.com",
    "date_of_birth": "1990-01-02T03:04:05Z"
  }'
```

Get a person:

```sh
curl -i http://localhost:8080/3f93df6d-ff51-4740-9d27-fc6b2f30281c
```

## Tests

Run all tests:

```sh
go test ./...
```

The integration test in `tests/` runs the whole application process with `go run ./cmd/api`, a temporary SQLite database, and a free local HTTP port. It does not use mocks or `httptest.NewServer`.

## Notes

SQLite keeps the service simple and easy to run locally or in a small container. For higher-concurrency production workloads, PostgreSQL would usually be preferred.
