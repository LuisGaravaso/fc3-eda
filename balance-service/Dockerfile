FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y librdkafka-dev

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -o balancecore ./cmd/balancecore/main.go

FROM golang:1.24

RUN apt-get update && apt-get install -y librdkafka-dev curl

WORKDIR /app
COPY --from=builder /app/balancecore .

EXPOSE 8080
CMD ["./balancecore"]