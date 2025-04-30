#!/bin/sh
set -e

: "${ACCESS_TOKEN_SECRET:=$(openssl rand -hex 32)}"
: "${REFRESH_TOKEN_SECRET:=$(openssl rand -hex 32)}"

export ACCESS_TOKEN_SECRET REFRESH_TOKEN_SECRET

exec /app/auth-api "$@"
