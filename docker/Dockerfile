# ---- Build Stage ----
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-api ./cmd/api

# ---- Run Stage ----
FROM scratch AS final
WORKDIR /app
COPY --from=builder /app/auth-api ./auth-api
COPY --from=builder /app/.env.example .env.example
ENTRYPOINT ["/app/auth-api"]
