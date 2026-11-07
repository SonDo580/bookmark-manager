# Bookmark manager

## Development

- Create `.env` file:
  - Check `internal/config/config.go`.

- First-time installation:

```bash
go mod tidy
```

- Run the application:

```bash
go run cmd/api/main.go
```

- Install extra dependencies:

```bash
go get ...
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

- Run all tests:

```bash
go test -v ./...
```
