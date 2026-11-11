FROM golang:1.27.3-alpine AS base

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go build -o bookmark_service cmd/api/main.go

FROM alpine:3.23.2 AS run

WORKDIR /app

COPY --from=base /opt/app/bookmark_service .

CMD ["/app/bookmark_service"]