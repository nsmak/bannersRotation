package app

import (
	"context"

	"github.com/nsmak/bannersRotation/internal/utils"
)

type domainError struct {
	BaseError
}

func newError(msg string, err error) *domainError {
	return &domainError{BaseError: BaseError{Message: msg, Err: err}}
}

// RotatorDomain отвечает за работу с баннерами.
type RotatorDomain struct {
	store Storage
	log   Logger
}

// NewRotator - возвращает новый инстанс домена.
func NewRotator(s Storage, l Logger) *RotatorDomain {
	return &RotatorDomain{store: s, log: l}
}

// AddBannerToSlot - добавляет новый баннер в ротацию в данном слоте.
func (r *RotatorDomain) AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error {
	err := r.store.AddBannerToSlot(ctx, bannerID, slotID)
	if err != nil {
		r.log.Error("can't add banner to slot", r.log.String("msg", err.Error()))
		return newError("add banner to slot error", err)
	}

	return nil
}

// RemoveBannerFromSlot - удаляет баннер из ротации в данном слоте.
func (r *RotatorDomain) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error {
	err := r.store.RemoveBannerFromSlot(ctx, bannerID, slotID)
	if err != nil {
		r.log.Error("can't remove banner from slot", r.log.String("msg", err.Error()))
		return newError("remove banner from slot error", err)
	}

	return nil
}

func (r *RotatorDomain) BannerIDForSlot(ctx context.Context, slotID, socialID int64) (int64, error) {
	stats, err := r.store.BannersStatistics(ctx, slotID, socialID)
	if err != nil {
		r.log.Error("can't get statistics about slot", r.log.String("msg", err.Error()))
		return 0, newError("slot statistics error", err)
	}

	lenStats := len(stats)

	showsCount := make([]int64, lenStats)
	clicksCount := make([]int64, lenStats)

	for i, s := range stats {
		showsCount[i] = s.ShowCount
		clicksCount[i] = s.ClickCount.Int64
	}

	index, err := utils.PlayWithBandit(showsCount, clicksCount)
	if err != nil {
		r.log.Error("play with bandit error", r.log.String("msg", err.Error()))
		return 0, newError("play with bandit error", err)
	}

	bannerID := stats[index].BannerID

	err = r.store.AddViewForBanner(ctx, bannerID, slotID, socialID)
	if err != nil {
		r.log.Error("add view for banner error", r.log.String("msg", err.Error()))
		return 0, newError("add view for banner error", err)
	}

	return bannerID, nil
}

func (r *RotatorDomain) AddClickForBanner(ctx context.Context, bannerID, slotID, socialID int64) error {
	err := r.store.AddClickForBanner(ctx, bannerID, slotID, socialID)
	if err != nil {
		r.log.Error("add click for banner error", r.log.String("msg", err.Error()))
		return newError("add click for banner error", err)
	}

	return nil
}
