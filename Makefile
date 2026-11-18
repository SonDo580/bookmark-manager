.PHONY: swagger run dev test docker-build docker-test docker-login docker-push

swagger:
	swag init -g cmd/api/main.go

run:
	go run cmd/api/main.go

dev: swagger run

COVERAGE_EXCLUDE := mocks|main.go|test|docs
COVERAGE_THRESHOLD := 80

test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html

	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

COVERAGE_FOLDER := ./test-output

docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --build-arg COVERAGE_EXCLUDE="$(COVERAGE_EXCLUDE)" --progress=plain --target test-out -t test:test --output $(COVERAGE_FOLDER) .

	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

IMG_NAME := sondo580/bookmark-service
IMG_TAG := latest
GIT_TAG = $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)

ifneq ($(GIT_TAG),)
	IMG_TAG := $(GIT_TAG)
endif

docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .

DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

docker-push:
	docker push $(IMG_NAME):$(IMG_TAG)