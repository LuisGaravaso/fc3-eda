FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y librdkafka-dev

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -o walletcore ./cmd/walletcore/main.go

FROM golang:1.24

RUN apt-get update && apt-get install -y librdkafka-dev curl

WORKDIR /app
COPY --from=builder /app/walletcore .

EXPOSE 8080
CMD ["./walletcore"]