FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN apk add build-base && go mod download

COPY . ./

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /gowatch


FROM alpine AS base

WORKDIR /app

COPY --from=builder /gowatch ./gowatch

RUN mkdir /config
ENV GOWATCH_DATABASE_DSN "/config/database.db"

ENTRYPOINT ["./gowatch"]