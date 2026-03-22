# Stage 1: Build
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o spaincoin ./node/cmd/

# Stage 2: Run (minimal image)
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/spaincoin .

# Data directory for blockchain storage
VOLUME ["/data"]

# P2P port
EXPOSE 30303
# RPC API port
EXPOSE 8545

CMD ["./spaincoin"]
