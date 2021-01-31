package app

import "database/sql"

type SocialGroup struct {
	ID          int64  `db:"id"`
	Description string `db:"description"`
}

type BannerSummary struct {
	BannerID   int64         `db:"banner_id"`
	SlotID     int64         `db:"slot_id"`
	SocialID   int64         `db:"social_id"`
	ShowCount  int64         `db:"show_count"`
	ClickCount sql.NullInt64 `db:"click_count"`
}

type BannerStatistic struct {
	BannerID int64   `db:"banner_id"`
	SlotID   int64   `db:"slot_id"`
	SocialID int64   `db:"social_id"`
	Date     float64 `db:"date"`
}

type StatType string

type MQBannerStatistic struct {
	Type     StatType
	BannerID int64
	SlotID   int64
	SocialID int64
	Date     float64
}

func NewMQBannerStatistic(statType StatType, stat BannerStatistic) MQBannerStatistic {
	return MQBannerStatistic{
		Type:     statType,
		BannerID: stat.BannerID,
		SlotID:   stat.SlotID,
		SocialID: stat.SocialID,
		Date:     stat.Date,
	}
}
