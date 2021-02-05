package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nsmak/bannersRotation/internal/app"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var errStore = errors.New("store error")

type RotatorDomainSuite struct {
	suite.Suite
	mockCtl   *gomock.Controller
	mockStore *MockStorage
	rotator   *app.RotatorDomain
}

func (s *RotatorDomainSuite) SetupTest() {
	s.mockCtl = gomock.NewController(s.T())
	s.mockStore = NewMockStorage(s.mockCtl)
	s.rotator = app.NewRotator(s.mockStore, &mockLogger{})
}

func (s *RotatorDomainSuite) TearDownTest() {
	s.mockCtl.Finish()
}

func (s *RotatorDomainSuite) TestAddBannerToSlotSuccess() {
	var bannerID int64 = 1
	var slotID int64 = 2
	ctx := context.Background()

	s.mockStore.EXPECT().AddBannerToSlot(ctx, bannerID, slotID).Return(nil)
	err := s.rotator.AddBannerToSlot(ctx, bannerID, slotID)

	s.Require().NoError(err)
}

func (s *RotatorDomainSuite) TestAddBannerToSlotFail() {
	var bannerID int64 = 1
	var slotID int64 = 2
	ctx := context.Background()

	s.mockStore.EXPECT().AddBannerToSlot(ctx, bannerID, slotID).Return(errStore)
	err := s.rotator.AddBannerToSlot(ctx, bannerID, slotID)

	s.Require().Error(err)
}

func (s *RotatorDomainSuite) TestRemoveBannerFromSlotSuccess() {
	var bannerID int64 = 1
	var slotID int64 = 2
	ctx := context.Background()

	s.mockStore.EXPECT().RemoveBannerFromSlot(ctx, bannerID, slotID).Return(nil)
	err := s.rotator.RemoveBannerFromSlot(ctx, bannerID, slotID)

	s.Require().NoError(err)
}

func (s *RotatorDomainSuite) TestRemoveBannerFromSlotFail() {
	var bannerID int64 = 1
	var slotID int64 = 2
	ctx := context.Background()

	s.mockStore.EXPECT().RemoveBannerFromSlot(ctx, bannerID, slotID).Return(errStore)
	err := s.rotator.RemoveBannerFromSlot(ctx, bannerID, slotID)

	s.Require().Error(err)
}

func (s *RotatorDomainSuite) TestBannerIDForSlotSuccess() {
	var slotID int64 = 1
	var socialID int64 = 1
	stats := mockStatistics()
	var expected int64 = 3
	ctx := context.Background()

	s.mockStore.EXPECT().BannersStatistics(ctx, slotID, socialID).Return(stats, nil)
	s.mockStore.EXPECT().AddViewForBanner(ctx, int64(3), slotID, socialID).Return(nil)
	bannerID, err := s.rotator.BannerIDForSlot(ctx, slotID, socialID)

	s.Require().NoError(err)
	s.Require().Equal(expected, bannerID)
}

func (s *RotatorDomainSuite) TestBannerIDForSlotStoreStatsFail() {
	var slotID int64 = 1
	var socialID int64 = 1
	ctx := context.Background()

	s.mockStore.EXPECT().BannersStatistics(ctx, slotID, socialID).Return(nil, errStore)
	bannerID, err := s.rotator.BannerIDForSlot(ctx, slotID, socialID)

	s.Require().Error(err)
	s.Require().True(errors.Is(err, errStore))
	s.Equal(int64(0), bannerID)
}

func (s *RotatorDomainSuite) TestBannerIDForSlotStoreAddViewFail() {
	var slotID int64 = 1
	var socialID int64 = 1
	stats := mockStatistics()
	ctx := context.Background()

	s.mockStore.EXPECT().BannersStatistics(ctx, slotID, socialID).Return(stats, nil)
	s.mockStore.EXPECT().AddViewForBanner(ctx, int64(3), slotID, socialID).Return(errStore)
	bannerID, err := s.rotator.BannerIDForSlot(ctx, slotID, socialID)

	s.Require().Error(err)
	s.Require().True(errors.Is(err, errStore))
	s.Equal(int64(0), bannerID)
}

func (s *RotatorDomainSuite) TestAddClickForBannerSuccess() {
	var bannerID int64 = 1
	var slotID int64 = 1
	var socialID int64 = 1
	ctx := context.Background()

	s.mockStore.EXPECT().AddClickForBanner(ctx, bannerID, slotID, socialID).Return(nil)
	err := s.rotator.AddClickForBanner(ctx, bannerID, slotID, socialID)

	s.Require().NoError(err)
}

func (s *RotatorDomainSuite) TestAddClickForBannerFail() {
	var bannerID int64 = 1
	var slotID int64 = 1
	var socialID int64 = 1
	ctx := context.Background()

	s.mockStore.EXPECT().AddClickForBanner(ctx, bannerID, slotID, socialID).Return(errStore)
	err := s.rotator.AddClickForBanner(ctx, bannerID, slotID, socialID)

	s.Require().Error(err)
	s.Require().True(errors.Is(err, errStore))
}

func TestRotatorDomainSuite(t *testing.T) {
	suite.Run(t, new(RotatorDomainSuite))
}

func mockStatistics() []app.BannerSummary {
	return []app.BannerSummary{
		{
			BannerID:   1,
			SlotID:     1,
			SocialID:   1,
			ShowCount:  6,
			ClickCount: 1,
		},
		{
			BannerID:   2,
			SlotID:     1,
			SocialID:   1,
			ShowCount:  7,
			ClickCount: 2,
		},
		{
			BannerID:   3,
			SlotID:     1,
			SocialID:   1,
			ShowCount:  5,
			ClickCount: 1,
		},
	}
}

type mockLogger struct {
}

func (m *mockLogger) Info(msg string, fields ...zap.Field) {
}

func (m *mockLogger) Warn(msg string, fields ...zap.Field) {
}

func (m *mockLogger) Error(msg string, fields ...zap.Field) {
}

func (m *mockLogger) String(key string, val string) zap.Field {
	return zap.Field{}
}

func (m *mockLogger) Int64(key string, val int64) zap.Field {
	return zap.Field{}
}

func (m *mockLogger) Duration(key string, val time.Duration) zap.Field {
	return zap.Field{}
}
