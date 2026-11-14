# Bookmark manager

## Development

- Create `.env` file: check `.env.example`.

- First-time installation:

```bash
go mod tidy
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

# generate docs + run
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

## Redis

- Use Redis CLI:

```bash
docker exec -it bookmark-service-redis redis-cli
```

## Deploy

- Build the image:

```bash
# use default tag (latest)
make docker-build
# OR with specific version tag
make docker-build TAG=v1.0.0
```

- Login to Docker Hub:

```bash
docker login
```

- Push the image to Docker Hub:

```bash
# use default tag (latest)
make docker-push
# OR with specific version tag
make docker-push TAG=v1.0.0
```

- Build and push in 1 go _(make sure to login first)_:

```bash
# use default tag (latest)
make docker-release
# OR with specific version tag
make docker-release TAG=v1.0.0
```
