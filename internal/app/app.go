package app

	"context"
type Logger interface {
}

type Storage interface {
	AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error
	RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error
	BannersStatistics(ctx context.Context, slotID, socialID int64) ([]BannerSummary, error)
	AddViewForBanner(ctx context.Context, bannerID, slotID, socialID int64) error
	AddClickForBanner(ctx context.Context, bannerID, slotID, socialID int64) error
	BannersShowStatisticsFilterByDate(ctx context.Context, from int64, to int64) ([]BannerStatistic, error)
	BannersClickStatisticsFilterByDate(ctx context.Context, from int64, to int64) ([]BannerStatistic, error)
}
