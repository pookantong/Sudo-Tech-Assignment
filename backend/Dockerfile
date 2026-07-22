# --- build stage ---
FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download 2>/dev/null || true

COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/api

# --- run stage ---
FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server .

EXPOSE 8080
ENTRYPOINT ["./server"]