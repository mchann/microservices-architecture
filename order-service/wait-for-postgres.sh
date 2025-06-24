#!/bin/sh

echo "Menunggu Postgres siap..."

until pg_isready -h $POSTGRES_HOST -U $POSTGRES_USER; do
  echo "Postgres belum siap - tunggu dulu..."
  sleep 2
done

echo "Postgres siap - menjalankan service"
exec ./order-service
