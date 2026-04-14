FROM golang:1.25.4-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o ./gateway-service ./internal/gateway-service/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /build/gateway-service ./gateway-service

ENTRYPOINT ["/app/gateway-service"]
CMD ["application"]