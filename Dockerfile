FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN apk add build-base && go mod download

COPY ./models ./models
COPY ./notifiers ./notifiers
COPY ./web ./web
COPY ./main.go ./main.go

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /gowatch


FROM alpine AS base

WORKDIR /app

COPY --from=builder /gowatch /app/gowatch
RUN mkdir /config

RUN addgroup -S gowatch && \
    adduser -S gowatch -G gowatch && \
    chown gowatch:gowatch /app && \
    chown gowatch:gowatch /config

USER gowatch

ENV GOWATCH_DATABASE_DSN "/config/database.db"

ENTRYPOINT ["/app/gowatch"]