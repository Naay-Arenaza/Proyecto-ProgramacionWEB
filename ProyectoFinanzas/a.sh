#!/bin/bash

CONTAINER_NAME=contenedor_ProyectoFinanzas
echo "evantando contenedor de Postgres..."
sudo docker compose up -d database

echo "Esperando a que Postgres esté disponible..."
until docker exec "$CONTAINER_NAME" bash -c "cat /dev/null > /dev/tcp/localhost/5432" 2>/dev/null; do
  echo "⏳ Esperando a que Postgres abra el puerto 5432..."
  sleep 1
done

echo "Postgres está listo!"
go run main.go


