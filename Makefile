BIN_ROT := "./bin/rotator"
BIN_STAT := "./bin/statistic"

build:
	go build -v -o $(BIN_ROT) ./cmd/rotator

build-statistic:
	go build -v -o $(BIN_STAT) ./cmd/statistic

run:
	sh ./deployments/deploy.sh run

stop:
	sh ./deployments/deploy.sh stop

test:
	sh ./deployments/deploy.sh test