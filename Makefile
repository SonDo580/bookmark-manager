.PHONY: swagger run dev test docker-build docker-push docker-release

swagger:
	swag init -g cmd/api/main.go

run:
	go run cmd/api/main.go

dev: swagger run

COVERAGE_EXCLUDE=mocks|main.go|test

test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html

IMAGE_NAME ?= sondo580/bookmark-service
TAG ?= latest

docker-build: swagger
	docker build -t $(IMAGE_NAME):$(TAG) .

docker-push:
	docker push $(IMAGE_NAME):$(TAG)

docker-release: docker-build docker-push