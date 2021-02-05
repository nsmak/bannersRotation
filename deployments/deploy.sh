#!/bin/sh

DOCKER_VOLUME=deployments_dbdata

case "$1" in
run)
  echo Deploying project
  docker-compose -f deployments/docker-compose.yaml up -d --build
  ;;

stop)
  echo Stopping project
  docker-compose -f deployments/docker-compose.yaml down

  echo Removing docker volume $DOCKER_VOLUME
  docker volume rm $DOCKER_VOLUME
  ;;

test)
  	echo Deploying project
  	docker-compose -f deployments/docker-compose.yaml up -d --build
  	sleep 10

  	echo Deploying tests
  	docker-compose -f deployments/docker-compose-tests.yaml up --build
  	rc=$?

  	echo Stopping project and tests
  	docker-compose -f deployments/docker-compose.yaml -f deployments/docker-compose-tests.yaml down

  	echo Removing docker volume $DOCKER_VOLUME
    docker volume rm $DOCKER_VOLUME
  	 exit $rc
  	 ;;
*)

  echo Usage ./deploy.sh run|stop|test
  ;;
esac