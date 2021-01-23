package app

import (
	"context"
	"log"
)

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
		// r.log.Error("can't add banner to slot", r.log.String("msg", err.Error())) TODO: - сделать и убрать стд
		log.Println("can't add banner to slot", err.Error())
		return err
	}

	return nil
}

// RemoveBannerFromSlot - удаляет баннер из ротации в данном слоте.
func (r *RotatorDomain) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error {
	err := r.store.RemoveBannerFromSlot(ctx, bannerID, slotID)
	if err != nil {
		// r.log.Error("can't remove banner from slot", r.log.String("msg", err.Error())) TODO: - сделать и убрать стд
		log.Println("can't remove banner from slot", err.Error())
		return err
	}

	return nil
}
