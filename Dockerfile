FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0 generate

RUN CGO_ENABLED=0 go build -o bot ./cmd/bot

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bot .

CMD ["./bot"]
