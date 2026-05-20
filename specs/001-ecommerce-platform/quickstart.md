# Quickstart: E-Commerce Backend Platform

## Prerequisites

- Docker and Docker Compose installed
- Go toolchain matching `go.mod` (Go 1.26)
- `.env` file configured (or use `.env.example` as a baseline)

## Run the stack

1. Start dependencies (Postgres, RabbitMQ, NATS):

   ```bash
   make compose-up
   ```

2. Apply database migrations:

   ```bash
   make migrate-up
   ```

3. Run the REST API service:

   ```bash
   make run
   ```

4. Access API documentation (Swagger if enabled):

   - http://localhost:8080/swagger/index.html

## Run tests

- Unit tests:

  ```bash
  make test
  ```

- Integration tests:

  ```bash
  make integration-test
  ```

## Example request

```bash
curl -X GET "http://localhost:8080/api/v1/public/categories"
```

## Notes

- REST is the only supported protocol for this feature.
- Keep changes within existing modular monolith layers (`entity`, `usecase`, `repo`, `controller`).
