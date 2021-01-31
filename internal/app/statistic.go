package app

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

var (
	typeShow  StatType = "show"
	typeClick StatType = "click"
)

type Statistic struct {
	log      Logger
	storage  Storage
	producer MQProducer
	interval time.Duration
}

func NewStatistic(logger Logger, storage Storage, producer MQProducer, interval time.Duration) *Statistic {
	return &Statistic{log: logger, storage: storage, producer: producer, interval: interval}
}

func (s *Statistic) Run(ctx context.Context) {
	doneCh := make(chan struct{})
	go startWorker(ctx, doneCh, s.interval, func() {
		s.publishStatisticMessage(ctx)
	})
	<-doneCh
}

func (s *Statistic) publishStatisticMessage(ctx context.Context) {
	err := s.producer.OpenChannel()
	if err != nil {
		s.log.Error("can't open channel", s.log.String("msg", err.Error()))
		return
	}
	defer func() {
		err := s.producer.CloseChannel()
		if err != nil {
			s.log.Error("can't close channel", s.log.String("msg", err.Error()))
		}
	}()

	from := time.Now()
	to := from.Add(s.interval)

	shows, err := s.storage.BannersShowStatisticsFilterByDate(ctx, from.Unix(), to.Unix())
	if err != nil {
		log.Println(err.Error())
		// s.log.Error("can't get events", s.log.String("msg", err.Error()))
		return
	}

	clicks, err := s.storage.BannersClickStatisticsFilterByDate(ctx, from.Unix(), to.Unix())
	if err != nil {
		log.Println(err.Error())
		// s.log.Error("can't get events", s.log.String("msg", err.Error()))
		return
	}

	for _, show := range shows {
		data, err := json.Marshal(NewMQBannerStatistic(typeShow, show))
		if err != nil {
			s.log.Error("can't marshal event notification", s.log.String("msg", err.Error()))
			continue
		}
		err = s.producer.Publish(data)
		if err != nil {
			s.log.Error("can't publish event notification", s.log.String("msg", err.Error()))
		}
	}

	for _, click := range clicks {
		data, err := json.Marshal(NewMQBannerStatistic(typeClick, click))
		if err != nil {
			s.log.Error("can't marshal event notification", s.log.String("msg", err.Error()))
			continue
		}
		err = s.producer.Publish(data)
		if err != nil {
			s.log.Error("can't publish event notification", s.log.String("msg", err.Error()))
		}
	}
}

func startWorker(ctx context.Context, done chan struct{}, interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			close(done)
			return
		case <-ticker.C:
			fn()
		}
	}
}
