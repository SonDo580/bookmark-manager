.PHONY: swagger run dev test

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
