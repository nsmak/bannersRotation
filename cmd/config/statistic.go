package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Rabbit struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	ExchangeName string `json:"exchange_name"`
	ExchangeType string `json:"exchange_type"`
	QueueName    string `json:"queue_name"`
	RoutingKey   string `json:"routing_key"`
	ConsumerTag  string `json:"consumer_tag"`
}

type Statistic struct {
	Logger        LoggerConf `json:"logger"`
	RabbitMQ      Rabbit     `json:"rabbit_mq"`
	Database      DBConf     `json:"database"`
	IntervalInSec int64      `json:"interval_in_sec"`
}

func NewStatistic(filePath string) (Statistic, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Statistic{}, fmt.Errorf("can't open config file: %w", err)
	}
	defer file.Close()

	var config Statistic
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Statistic{}, fmt.Errorf("can't decode config: %w", err)
	}
	return config, nil
}
