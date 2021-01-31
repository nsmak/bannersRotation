// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/nsmak/bannersRotation/internal/app"
	serverapi "github.com/nsmak/bannersRotation/internal/server/rest/api"
	sqlstorage "github.com/nsmak/bannersRotation/internal/storage/sql"
	"github.com/stretchr/testify/suite"
)

var (
	psUsr  = "postgres"
	psPass = "password"
	psAddr = "db:5432"
	psDB   = "postgres"

	restURL = "http://rotator:8888"
)

type IntegrationSuite struct {
	suite.Suite
	db      *sqlx.DB
	storage *sqlstorage.BannerDataStore
	slot    app.Slot
	banners []app.Banner
	group   app.SocialGroup
}

func (s *IntegrationSuite) SetupTest() {
	storage, err := sqlstorage.New(context.Background(), psUsr, psPass, psAddr, psDB)
	if err != nil {
		log.Fatalln(err.Error())
	}
	s.storage = storage
	s.db = s.initDB()
	s.slot = s.defaultSlot()
	s.banners = s.defaultBanners()
	s.group = s.defaultGroup()
	s.saveDefaultSlots()
	s.saveDefaultBanners()
	s.saveDefaultGroups()
	s.initBanners()
}

func (s *IntegrationSuite) TearDownTest() {
	s.removeStatisticsForDefaultData()
	s.removeDefaultBannersFromSlots()
	s.removeDefaultBanners()
	s.removeDefaultSlots()
	s.removeDefaultGroups()
	_ = s.db.Close()
	_ = s.storage.Close()
}

func (s *IntegrationSuite) initDB() *sqlx.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", psUsr, psPass, psAddr, psDB)
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	return db
}

func (s *IntegrationSuite) TestAllBannersWillShow() {
	for i := 0; i < 11; i++ {
		resp, err := http.Get(restURL + fmt.Sprintf("/banner?slot_id=%d&soc_dem_id=%d", s.slot.ID, s.group.ID))

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
	}

	stats, err := s.storage.BannersStatistics(context.Background(), s.slot.ID, s.group.ID)

	s.Require().NoError(err, "lol")

	for _, stat := range stats {
		s.Require().True(stat.ShowCount > 0, "banner id:", stat.BannerID)
	}
}

func (s *IntegrationSuite) TestClickestBannerWillGetMoreShows() {
	for i := 0; i < 11; i++ {
		resp, err := http.Get(restURL + fmt.Sprintf("/banner?slot_id=%d&soc_dem_id=%d", s.slot.ID, s.group.ID))

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
	}

	for i := 0; i < 6; i++ {
		form := serverapi.BannerClickFrom{BannerID: s.banners[0].ID, SlotID: s.slot.ID, SocDemID: s.group.ID}
		data, err := json.Marshal(&form)

		s.Require().NoError(err)

		resp, err := http.Post(restURL+"/banner/click/add", "application/json", bytes.NewReader(data))

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
	}

	for i := 0; i < 21; i++ {
		resp, err := http.Get(restURL + fmt.Sprintf("/banner?slot_id=%d&soc_dem_id=%d", s.slot.ID, s.group.ID))

		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
	}

	stats, err := s.storage.BannersStatistics(context.Background(), s.slot.ID, s.group.ID)

	s.Require().NoError(err)

	var bannerID int64
	var currentMaxVal int64

	for _, stat := range stats {
		if stat.ShowCount > currentMaxVal {
			currentMaxVal = stat.ShowCount
			bannerID = stat.BannerID
		}
		log.Println("banner id:", bannerID, "show count: ", stat.ShowCount)
	}
	log.Println("Winner banner id:", bannerID)

	s.Require().Equal(s.banners[0].ID, bannerID)
}

func (s *IntegrationSuite) defaultSlot() app.Slot {
	return app.Slot{ID: 111, Description: "Slot 111"}
}

func (s *IntegrationSuite) defaultBanners() []app.Banner {
	return []app.Banner{
		{
			ID:          100,
			Description: "Banner 100",
		},
		{
			ID:          101,
			Description: "Banner 101",
		},
		{
			ID:          102,
			Description: "Banner 102",
		},
	}
}

func (s *IntegrationSuite) defaultGroup() app.SocialGroup {
	return app.SocialGroup{ID: 105, Description: "SocialGroup 105"}
}

func (s *IntegrationSuite) saveDefaultSlots() {
	_, err := s.db.Exec("INSERT INTO slot (id, description) VALUES ($1, $2)", s.slot.ID, s.slot.Description)
	if err != nil {
		log.Fatalln(err.Error(), "save slot")
	}
}

func (s *IntegrationSuite) saveDefaultBanners() {
	for _, banner := range s.defaultBanners() {
		_, err := s.db.Exec("INSERT INTO banner (id, description) VALUES ($1, $2)", banner.ID, banner.Description)
		if err != nil {
			log.Fatalln(err.Error(), "save banners")
		}
	}
}

func (s *IntegrationSuite) saveDefaultGroups() {
	_, err := s.db.Exec("INSERT INTO social_dem (id, description) VALUES ($1, $2)", s.group.ID, s.group.Description)
	if err != nil {
		log.Fatalln(err.Error(), "save group")
	}
}

func (s *IntegrationSuite) initBanners() {
	for _, banner := range s.banners {
		err := s.storage.AddBannerToSlot(context.Background(), banner.ID, s.slot.ID)
		if err != nil {
			log.Fatalln(err.Error(), "init banners")
		}
	}
}

func (s *IntegrationSuite) removeDefaultBannersFromSlots() {
	_, err := s.db.Exec("DELETE FROM banner_slot WHERE slot_id=$1", s.slot.ID)
	if err != nil {
		log.Fatalln(err.Error(), "delete banner form slot")
	}
}

func (s *IntegrationSuite) removeStatisticsForDefaultData() {
	_, err := s.db.Exec("DELETE FROM banner_showing WHERE slot_id=$1", s.slot.ID)
	if err != nil {
		log.Fatalln(err.Error(), "delete shows stats")
	}

	_, err = s.db.Exec("DELETE FROM banner_click WHERE slot_id=$1", s.slot.ID)
	if err != nil {
		log.Fatalln(err.Error(), "delete click statistics")
	}
}

func (s *IntegrationSuite) removeDefaultSlots() {
	_, err := s.db.Exec("DELETE FROM slot WHERE id=$1", s.slot.ID)
	if err != nil {
		log.Fatalln(err.Error(), "delete slot")
	}
}

func (s *IntegrationSuite) removeDefaultBanners() {
	for _, banner := range s.defaultBanners() {
		_, err := s.db.Exec("DELETE FROM banner WHERE id=$1", banner.ID)
		if err != nil {
			log.Fatalln(err.Error(), "delete banners")
		}
	}
}

func (s *IntegrationSuite) removeDefaultGroups() {
	_, err := s.db.Exec("DELETE FROM social_dem WHERE id=$1", s.group.ID)
	if err != nil {
		log.Fatalln(err.Error(), "delete group")
	}
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationSuite))
}
