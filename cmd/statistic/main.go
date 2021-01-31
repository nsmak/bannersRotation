package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nsmak/bannersRotation/cmd/config"
	"github.com/nsmak/bannersRotation/internal/app"
	"github.com/nsmak/bannersRotation/internal/logger"
	"github.com/nsmak/bannersRotation/internal/mq/rabbit"
	sqlstorage "github.com/nsmak/bannersRotation/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/statistic.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.NewStatistic(configFile)
	if err != nil {
		log.Fatalf("can't get config: %v", err)
	}

	logg, err := logger.New(cfg.Logger.Level, cfg.Logger.FilePath)
	if err != nil {
		log.Fatalf("can't start logger %v\n", err)
	}

	producer, err := rabbit.NewProducer(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("can't create consumer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("starting store service")
	storage, err := sqlstorage.New(
		ctx,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Address,
		cfg.Database.DBName,
	)
	if err != nil {
		log.Fatalf("failed to start storage connection: " + err.Error()) // nolint: gocritic
	}

	statistic := app.NewStatistic(logg, storage, producer, time.Duration(cfg.IntervalInSec)*time.Second)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		<-signals
		signal.Stop(signals)
		err := producer.CloseConn()
		if err != nil {
			logg.Error("can't close connection", logg.String("msg", err.Error()))
		}
		cancel()
	}()

	log.Println("starting statistic service")
	statistic.Run(ctx)
}
