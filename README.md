# Bookmark manager

## Development

- Create `.env` file: check `.env.example`.

- First-time installation:

```bash
go mod tidy
```

- Generate key pair for authentication:

```bash
make generate-key
```

- Start services:

```bash
docker compose up -d
```

- Run the application in development mode:

```bash
# generate docs
make swagger

# run
make run

# (combined) generate docs + run
make dev
```

- Check API documentation at `http://localhost:8080/swagger/index.html` by default.

- Install extra dependencies / tools:

```bash
go get ...
go install ...
```

## Testing

- Mock creation example:

```go
//go:generate mockery --name=Health --output=./mocks --outpkg=mocks
type Health interface {
    Check() model.HealthCheckResponse
}

// Then run: go generate
```

- Run all tests (with coverage report):

```bash
make test
```

## Tools

- Interact with Redis: use Redis CLI.

```bash
docker exec -it bookmark-service-redis redis-cli
```

- Interact with DB: use DBeaver (or other tools).
