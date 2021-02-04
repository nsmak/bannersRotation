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
	"github.com/nsmak/bannersRotation/internal/server/rest"
	"github.com/nsmak/bannersRotation/internal/server/rest/api"
	sqlstorage "github.com/nsmak/bannersRotation/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/rotator.json", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.NewCalendar(configFile)
	if err != nil {
		log.Fatalf("can't get config: %v", err)
	}

	logg, err := logger.New(cfg.Logger.Level, cfg.Logger.FilePath)
	if err != nil {
		log.Fatalf("can't start logger %v\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Println("starting store service")
	storage, err := sqlstorage.New(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Address, cfg.DB.DBName)
	if err != nil {
		log.Fatalf("failed to start storage connection: " + err.Error()) // nolint: gocritic
	}

	rotator := app.NewRotator(storage, logg)
	server := rest.NewServer(api.New(rotator), cfg.RestServer.Address, logg)

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		<-signals
		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		log.Println("stopping rest server...")
		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop rest server", logg.String("msg", err.Error()))
		}
	}()

	log.Println("starting REST server at " + server.Address)
	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start rest server", logg.String("msg", err.Error()))
	}
	log.Println("server stopped")
}
