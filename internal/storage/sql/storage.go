package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib" // nolint: gci
	"github.com/jmoiron/sqlx"
	"github.com/nsmak/bannersRotation/internal/app"
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

func (s *BannerDataStore) Close() error {
	return s.db.Close()
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

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return storage.NewError("can't start transactions", err)
	}
	defer tx.Rollback() // nolint: errcheck

	_, err = tx.ExecContext(ctx, "INSERT INTO banner_slot (banner_id, slot_id) VALUES ($1, $2)", bannerID, slotID)
	if err != nil {
		return storage.NewError("can't add banner into slot", err)
	}

	groups, err := s.socialGroups(ctx)
	if err != nil {
		return storage.NewError("can't get social groups", err)
	}

	for _, grp := range groups {
		_, err := tx.ExecContext(
			ctx,
			"INSERT INTO banner_showing (banner_id, slot_id, social_id, date) VALUES ($1, $2, $3, current_timestamp)",
			bannerID, slotID, grp.ID,
		)
		if err != nil {
			return storage.NewError("can't add view for banner", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return storage.NewError("can't commit transactions", err)
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

func (s *BannerDataStore) BannersStatistics(ctx context.Context, slotID, socialID int64) ([]app.BannerSummary, error) {
	slotIsExist, err := s.slotIsExist(ctx, slotID)
	if err != nil {
		return nil, storage.NewError("can't get info about slot", err)
	}

	if !slotIsExist {
		return nil, storage.ErrSlotNotFound
	}

	socialIsExist, err := s.socialGroupIsExist(ctx, socialID)
	if err != nil {
		return nil, storage.NewError("can't get info about social group", err)
	}

	if !socialIsExist {
		return nil, storage.ErrSocialGroupNotFound
	}

	var stats []app.BannerSummary
	err = s.db.SelectContext(
		ctx,
		&stats,
		`SELECT sh.banner_id, sh.slot_id, sh.social_id, count(sh.*) show_count, cl.count click_count
			FROM banner_showing sh
			LEFT JOIN (SELECT banner_id, slot_id, social_id, count(date) FROM banner_click GROUP BY 1,2,3) cl
			ON (sh.slot_id=cl.slot_id AND sh.banner_id=cl.banner_id AND sh.social_id=cl.social_id)
			WHERE sh.slot_id=$1 AND sh.social_id=$2
			GROUP BY 1,2,3,5
		`,
		slotID, socialID,
	)
	if err != nil {
		return nil, storage.NewError("can't get statistics", err)
	}

	if len(stats) == 0 {
		return nil, storage.ErrStatisticsNotFound
	}

	return stats, nil
}

func (s *BannerDataStore) AddViewForBanner(ctx context.Context, bannerID, slotID, socialID int64) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO banner_showing (banner_id, slot_id, social_id, date) VALUES ($1, $2, $3, current_timestamp)",
		bannerID, slotID, socialID,
	)
	if err != nil {
		return storage.NewError("can't add view for banner", err)
	}

	return nil
}

func (s *BannerDataStore) AddClickForBanner(ctx context.Context, bannerID, slotID, socialID int64) error {
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

	socialIsExist, err := s.socialGroupIsExist(ctx, socialID)
	if err != nil {
		return storage.NewError("can't get info about social group", err)
	}

	if !socialIsExist {
		return storage.ErrSocialGroupNotFound
	}

	_, err = s.db.ExecContext(
		ctx,
		"INSERT INTO banner_click (banner_id, slot_id, social_id, date) VALUES ($1, $2, $3, current_timestamp)",
		bannerID, slotID, socialID,
	)
	if err != nil {
		return storage.NewError("can't add view for banner", err)
	}

	return nil
}

func (s *BannerDataStore) BannersShowStatisticsFilterByDate(ctx context.Context, from int64, to int64) ([]app.BannerStatistic, error) {
	var shows []app.BannerStatistic
	err := s.db.SelectContext(
		ctx,
		&shows,
		`SELECT banner_id, slot_id, social_id, extract(epoch from date) date 
			FROM banner_showing 
			WHERE extract(epoch from date) >=$1 AND extract(epoch from date) <=$2`,
		from, to,
	)
	if err != nil {
		return nil, storage.NewError("can't get shows info", err)
	}

	return shows, nil
}

func (s *BannerDataStore) BannersClickStatisticsFilterByDate(ctx context.Context, from int64, to int64) ([]app.BannerStatistic, error) {
	var shows []app.BannerStatistic
	err := s.db.SelectContext(
		ctx,
		&shows,
		`SELECT banner_id, slot_id, social_id, extract(epoch from date) date  
			FROM banner_click 
			WHERE extract(epoch from date) >=$1 AND extract(epoch from date) <=$2`,
		from, to,
	)
	if err != nil {
		return nil, storage.NewError("can't get shows info", err)
	}

	return shows, nil
}

func (s *BannerDataStore) socialGroups(ctx context.Context) ([]app.SocialGroup, error) {
	var groups []app.SocialGroup
	err := s.db.SelectContext(ctx, &groups, "SELECT id, description FROM social_dem")
	if err != nil {
		return nil, storage.NewError("can't get social groups", err)
	}

	return groups, nil
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

func (s *BannerDataStore) socialGroupIsExist(ctx context.Context, socialID int64) (bool, error) {
	var count int

	err := s.db.GetContext(ctx, &count, "SELECT  COUNT(*) FROM social_dem WHERE id=$1", socialID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, storage.NewError("can't get exist info about social group", err)
	}

	return count > 0, nil
}
