#!/bin/sh

cleanup() {
  echo "Cleaning up..."
  docker container prune -f
  docker image prune -f
  docker volume prune -f
}

# Set trap to call cleanup on exit or interrupt (Ctrl+C)
trap cleanup EXIT

docker-compose up
