#!/bin/sh

DOCKER_NETWORK=rotator_network
DOCKER_VOLUME=deployments_dbdata

case "$1" in
run)
  echo Creating network $DOCKER_NETWORK
  docker network create $DOCKER_NETWORK

  echo Deploying project
  docker-compose -f deployments/docker-compose.yaml up -d --build
  ;;

stop)
  echo Stopping project
  docker-compose -f deployments/docker-compose.yaml down

  echo Removing docker network $DOCKER_NETWORK
  docker network rm $DOCKER_NETWORK

  echo Removing docker volume $DOCKER_VOLUME
  docker volume rm $DOCKER_VOLUME
  ;;

*)
  echo Usage ./deploy.sh up|down
  ;;
esac