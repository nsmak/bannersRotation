package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib" // nolint: gci
	"github.com/jmoiron/sqlx"
	"github.com/nsmak/bannersRotation/internal/storage"
)

type BannerDataStore struct {
	db *sqlx.DB
}

func New(ctx context.Context, user, pass, addr, dbName string) (*BannerDataStore, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, addr, dbName)
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, storage.NewError("can't open db store", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, storage.NewError("ping error", err)
	}

	return &BannerDataStore{db: db}, nil
}

func (s *BannerDataStore) AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error {
	bannerInSlotExist, err := s.bannerIsExistInSlot(ctx, bannerID, slotID)
	if err != nil {
		return storage.NewError("can't get info about banner in slot", err)
	}

	if bannerInSlotExist {
		return storage.ErrBannerInSlotAlreadyExist
	}

	bannerIsExist, err := s.bannerIsExist(ctx, bannerID)
	if err != nil {
		return storage.NewError("can't get info about banner", err)
	}

	if !bannerIsExist {
		return storage.ErrBannerNotFound
	}

	slotIsExist, err := s.slotIsExist(ctx, slotID)
	if err != nil {
		return storage.NewError("can't get info about slot", err)
	}

	if !slotIsExist {
		return storage.ErrSlotNotFound
	}

	_, err = s.db.ExecContext(
		ctx,
		"INSERT INTO banner_slot (banner_id, slot_id) VALUES ($1, $2)", bannerID, slotID,
	)
	if err != nil {
		return storage.NewError("can't add banner into slot", err)
	}

	return nil
}

func (s *BannerDataStore) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error {
	slotIsExist, err := s.slotIsExist(ctx, slotID)
	if err != nil {
		return storage.NewError("can't get info about slot", err)
	}

	if !slotIsExist {
		return storage.ErrSlotNotFound
	}

	bannerIsExist, err := s.bannerIsExistInSlot(ctx, bannerID, slotID)
	if err != nil {
		return storage.NewError("can't get info about banner in slot", err)
	}

	if !bannerIsExist {
		return storage.ErrBannerInSlotNotFound
	}

	_, err = s.db.ExecContext(ctx, "DELETE FROM banner_slot WHERE banner_id=$1 AND slot_id=$2", bannerID, slotID)
	if err != nil {
		return storage.NewError("can't remove banner from slot", err)
	}

	return nil
}

func (s *BannerDataStore) bannerIsExist(ctx context.Context, bannerID int64) (bool, error) {
	var count int

	err := s.db.GetContext(ctx, &count, "SELECT  COUNT(*) FROM banner WHERE id=$1", bannerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, storage.NewError("can't get exist info about banner", err)
	}

	return count > 0, nil
}

func (s *BannerDataStore) slotIsExist(ctx context.Context, slotID int64) (bool, error) {
	var count int

	err := s.db.GetContext(ctx, &count, "SELECT  COUNT(*) FROM slot WHERE id=$1", slotID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, storage.NewError("can't get exist info about slot", err)
	}

	return count > 0, nil
}

func (s *BannerDataStore) bannerIsExistInSlot(ctx context.Context, bannerID, slotID int64) (bool, error) {
	var count int

	err := s.db.GetContext(
		ctx,
		&count, "SELECT COUNT(*) FROM banner_slot WHERE banner_id=$1 AND slot_id=$2", bannerID, slotID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, storage.NewError("can't get exist info about banner in slot", err)
	}

	return count > 0, nil
}
