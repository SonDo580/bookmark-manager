# Bookmark manager

## Development

- Create `.env` file: check `.env.example`.

- First-time installation:

```bash
go mod tidy
```

- Run the application in development mode:

```bash
make dev
```

- API documentation: at `http://localhost:8080/swagger/index.html` by default.

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

- Run all tests (with code coverage report):

```bash
make test
```

## Redis

- Use Redis CLI:

```bash
docker exec -it bookmark-service-redis redis-cli
```
