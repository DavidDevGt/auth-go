# ---------- build stage ----------
FROM --platform=linux/amd64 golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN [ ! -f go.mod ] && go mod init auth-api || true
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o auth-api .
# ---------- final stage ----------
FROM --platform=linux/amd64 alpine:latest
RUN apk add --no-cache openssl netcat-openbsd

WORKDIR /app
COPY --from=builder /app/auth-api .
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /app/auth-api /entrypoint.sh
ENV ACCESS_TOKEN_SECRET "" \
    REFRESH_TOKEN_SECRET ""
ENTRYPOINT ["/entrypoint.sh"]