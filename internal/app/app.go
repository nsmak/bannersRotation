package app

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	String(key string, val string) zap.Field
	Int64(key string, val int64) zap.Field
	Duration(key string, val time.Duration) zap.Field
}

//go:generate mockgen -destination=./mock_storage_test.go -package=app_test . Storage
//go:generate mockgen -destination=../server/rest/api/mock_storage_test.go -package=api_test . Storage
type Storage interface {
	AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error
	RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error
	BannersStatistics(ctx context.Context, slotID, socialID int64) ([]BannerSummary, error)
	AddViewForBanner(ctx context.Context, bannerID, slotID, socialID int64) error
	AddClickForBanner(ctx context.Context, bannerID, slotID, socialID int64) error
	BannersShowStatisticsFilterByDate(ctx context.Context, from int64, to int64) ([]BannerStatistic, error)
	BannersClickStatisticsFilterByDate(ctx context.Context, from int64, to int64) ([]BannerStatistic, error)
}

type MQProducer interface {
	Publish(body []byte) error
	OpenChannel() error
	CloseChannel() error
	CloseConn() error
}
