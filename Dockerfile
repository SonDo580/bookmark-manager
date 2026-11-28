# === Base ===

FROM golang:1.27.3-alpine AS base

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


# === Build ===

FROM base AS build

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6
RUN swag init -g cmd/api/main.go

RUN go build -o bookmark_service cmd/api/main.go


# === Test exec ===

FROM base AS test-exec

ARG _out_dir="/tmp/coverage"
ARG COVERAGE_EXCLUDE

RUN mkdir -p ${_out_dir} && \
    go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1 && \
	grep -vE "${COVERAGE_EXCLUDE}" coverage.tmp > ${_out_dir}/coverage.out && \
	go tool cover -html=${_out_dir}/coverage.out -o ${_out_dir}/coverage.html


# === Test output ===

FROM scratch AS test-out

ARG _out_dir="/tmp/coverage"

COPY --from=test-exec ${_out_dir}/coverage.out /
COPY --from=test-exec ${_out_dir}/coverage.html /


# === Run ===

FROM alpine:3.23.2 AS run

WORKDIR /app

COPY --from=build /opt/app/bookmark_service .

CMD ["/app/bookmark_service"]
