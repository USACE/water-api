#!/usr/bin/env bash

case $1 in

  start | up)
    echo -n "Starting minio..."
    docker compose -f docker-compose.minio.yml up --build
    ;;

  stop | down)
    echo -n "Stopping minio..."
    docker compose -f docker-compose.minio.yml down
    docker compose -f docker-compose.minio.yml down
    ;;

  *)
    echo -n "unknown arg"
    ;;
esac