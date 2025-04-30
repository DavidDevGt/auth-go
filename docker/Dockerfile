# ---------- build stage ----------
FROM --platform=linux/amd64 golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o auth-api .

# ---------- final stage ----------
FROM --platform=linux/amd64 alpine:latest
RUN apk add --no-cache openssl

WORKDIR /app
COPY --from=builder /app/auth-api .
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /app/auth-api /entrypoint.sh

ENV ACCESS_TOKEN_SECRET "" \
    REFRESH_TOKEN_SECRET ""

ENTRYPOINT ["/entrypoint.sh"]
