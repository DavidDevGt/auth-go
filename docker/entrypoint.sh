#!/bin/sh
set -e

: "${ACCESS_TOKEN_SECRET:=$(openssl rand -hex 32)}"
: "${REFRESH_TOKEN_SECRET:=$(openssl rand -hex 32)}"

export ACCESS_TOKEN_SECRET REFRESH_TOKEN_SECRET

# Esperar a que la base de datos esté lista
echo "Esperando para que la BD este lista :) Dei V"
MAX_RETRIES=30
COUNT=0
while [ $COUNT -lt $MAX_RETRIES ]; do
    nc -z db 3306 && break
    echo "La base de datos no está lista :c sad, reintentando ($COUNT/$MAX_RETRIES)..."
    COUNT=$((COUNT+1))
    sleep 1
done

if [ $COUNT -eq $MAX_RETRIES ]; then
    echo "Error: No se pudo conectar a la base de datos después de $MAX_RETRIES intentos. :c chale"
    exit 1
fi

echo "La BD esta lista, iniciando la app..."
exec /app/auth-api "$@"
