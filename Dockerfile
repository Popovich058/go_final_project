FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o taskscheduler main.go

ENTRYPOINT ["./taskscheduler"]
