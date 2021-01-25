#!/bin/sh

case "$1" in
up)
  docker network create rotator_network
  docker-compose -f deployments/docker-compose.yaml up -d --build
  ;;

down)
  docker-compose -f deployments/docker-compose.yaml down
  docker network rm rotator_network
  docker volume rm deployments_dbdata
  ;;

*)
  echo Usage ./deploy.sh up|down
  ;;
esac