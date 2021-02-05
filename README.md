# Banners Rotation [![Build Status](https://travis-ci.com/nsmak/bannersRotation.svg?branch=master)](https://travis-ci.com/nsmak/bannersRotation) [![Go Report Card](https://goreportcard.com/badge/github.com/nsmak/bannersRotation)](https://goreportcard.com/report/github.com/nsmak/bannersRotation)

Banner rotation service, based on  UCB1 algorithm (multiarmed bandit).

# Usage

## With docker
### Run
```
$ make run
```

### Stop
```
$ make stop
```

## Custom:
```
$ make build
```

### or

```
$ make build-statistic
```
### for build statistc sub service


## Sample rotator service config.json:

``` json 
{
  "logger": {
    "level": -1,
    "file_path": "./zap.log"
  },
  "rest_server": {
    "address": "rotator:8888"
  },
  "database": {
    "username": "postgres",
    "password": "password",
    "address": "db:5432",
    "db_name": "postgres"
  }
}
```

## Sample statistic service config.json:

``` json 
{
  "rabbit_mq": {
    "address": "mq:5672",
    "username": "guest",
    "password": "guest",
    "exchange_name": "stat_exchange",
    "exchange_type": "direct",
    "queue_name": "stat_queue",
    "routing_key": "stat_key",
    "consumer_tag": "stat_tag"
  },
  "logger": {
    "level": -1,
    "file_path": "./statistic.log"
  },
  "database": {
    "username": "postgres",
    "password": "password",
    "address": "db:5432",
    "db_name": "postgres"
  },
  "interval_in_sec": 60
}
```