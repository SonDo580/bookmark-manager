# === Build stage ===

FROM golang:1.27.3-alpine AS build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bookmark_service cmd/api/main.go

# === Run stage ===

FROM alpine:3.23.2 AS run

WORKDIR /app

COPY --from=build /opt/app/bookmark_service .

CMD ["/app/bookmark_service"]
