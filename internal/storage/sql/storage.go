package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib" // nolint: gci
	"github.com/jmoiron/sqlx"
	"github.com/nsmak/bannersRotation/internal/app"
	"github.com/nsmak/bannersRotation/internal/storage"
)

const (
	violatesForeignKeyConstraintCode = "23503"
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
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return storage.NewError("can't start transactions", err)
	}
	defer tx.Rollback() // nolint: errcheck

	_, err = tx.ExecContext(ctx, "INSERT INTO banner_slot (banner_id, slot_id) VALUES ($1, $2)", bannerID, slotID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == violatesForeignKeyConstraintCode {
				return storage.NewError(pgErr.Error(), storage.ErrObjectNotFound)
			}
		}

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
	_, err := s.db.ExecContext(ctx, "DELETE FROM banner_slot WHERE banner_id=$1 AND slot_id=$2", bannerID, slotID)
	if err != nil {
		return storage.NewError("can't remove banner from slot", err)
	}

	return nil
}

func (s *BannerDataStore) BannersStatistics(ctx context.Context, slotID, socialID int64) ([]app.BannerSummary, error) {
	rows, err := s.db.QueryxContext(
		ctx,
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
	defer rows.Close()

	var stats []app.BannerSummary
	for rows.Next() {
		var summary app.BannerSummary
		var clickCount sql.NullInt64

		err := rows.Scan(&summary.BannerID, &summary.SlotID, &summary.SocialID, &summary.ShowCount, &clickCount)
		if err != nil {
			return nil, storage.NewError("scan error", err)
		}
		summary.ClickCount = clickCount.Int64
		stats = append(stats, summary)
	}

	err = rows.Err()
	if err != nil {
		return nil, storage.NewError("rows error", err)
	}

	if len(stats) == 0 {
		return nil, storage.ErrObjectNotFound
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
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO banner_click (banner_id, slot_id, social_id, date) VALUES ($1, $2, $3, current_timestamp)",
		bannerID, slotID, socialID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == violatesForeignKeyConstraintCode {
				return storage.NewError(pgErr.Error(), storage.ErrObjectNotFound)
			}
		}

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
