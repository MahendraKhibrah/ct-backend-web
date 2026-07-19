#!/bin/sh
set -eu

: "${COMPOSE_PROJECT_NAME:?COMPOSE_PROJECT_NAME wajib diisi}"
: "${CONTAINER_NAME:?CONTAINER_NAME wajib diisi}"
: "${IMAGE_NAME:?IMAGE_NAME wajib diisi}"

DOCKER_NETWORK="${DOCKER_NETWORK:-web}"
STOP_TIMEOUT="${STOP_TIMEOUT:-20}"
DEPLOY_WAIT_TIMEOUT="${DEPLOY_WAIT_TIMEOUT:-90}"

if ! docker network inspect "$DOCKER_NETWORK" >/dev/null 2>&1; then
  docker network create "$DOCKER_NETWORK" >/dev/null 2>&1 \
    || docker network inspect "$DOCKER_NETWORK" >/dev/null
fi

echo "Memvalidasi konfigurasi Docker Compose..."
docker compose config --quiet

OLD_IMAGE_ID="$(docker image inspect "$IMAGE_NAME" --format '{{.Id}}' 2>/dev/null || true)"

echo "Membangun image $IMAGE_NAME..."
docker compose build

echo "Menghentikan deployment lama..."
docker compose down --remove-orphans --timeout "$STOP_TIMEOUT"
docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true

echo "Menjalankan deployment baru..."
docker compose up \
  --detach \
  --no-build \
  --pull never \
  --remove-orphans \
  --wait \
  --wait-timeout "$DEPLOY_WAIT_TIMEOUT"

docker compose ps

NEW_IMAGE_ID="$(docker image inspect "$IMAGE_NAME" --format '{{.Id}}')"
if [ -n "$OLD_IMAGE_ID" ] && [ "$OLD_IMAGE_ID" != "$NEW_IMAGE_ID" ]; then
  docker image rm "$OLD_IMAGE_ID" >/dev/null 2>&1 || true
fi
