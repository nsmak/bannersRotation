package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Rotator struct {
	Logger     LoggerConf `json:"logger"`
	RestServer RestConf   `json:"rest_server"`
	DB         DBConf     `json:"database"`
}

func NewCalendar(filePath string) (Rotator, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Rotator{}, fmt.Errorf("can't open config file: %w", err)
	}
	defer file.Close()

	var config Rotator
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Rotator{}, fmt.Errorf("can't decode config: %w", err)
	}
	return config, nil
}

type LoggerConf struct {
	Level    int8   `json:"level"`
	FilePath string `json:"file_path"`
}

type RestConf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type DBConf struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address  string `json:"address"`
	DBName   string `json:"db_name"`
}
