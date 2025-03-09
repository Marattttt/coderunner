FROM golang:1.24-alpine AS builder

WORKDIR /app/src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/server ./cmd/server/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY migrations/ .

ENTRYPOINT ["/app/server"]
