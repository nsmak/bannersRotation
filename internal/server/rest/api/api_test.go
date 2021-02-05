package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/nsmak/bannersRotation/internal/app"
	serverapi "github.com/nsmak/bannersRotation/internal/server/rest/api"
	"github.com/nsmak/bannersRotation/internal/storage"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

var errUnknown = errors.New("unknown")

type ApiSuite struct {
	suite.Suite
	mockCtl   *gomock.Controller
	mockStore *MockStorage
	server    *httptest.Server
	ctx       context.Context
	api       *serverapi.API
	rotator   *app.RotatorDomain
}

func (s *ApiSuite) SetupTest() {
	s.mockCtl = gomock.NewController(s.T())
	s.mockStore = NewMockStorage(s.mockCtl)
	s.ctx = context.Background()
	s.rotator = app.NewRotator(s.mockStore, &mockLogger{})
	s.api = serverapi.New(s.rotator)

	router := mux.NewRouter()
	for _, route := range s.api.Routes() {
		route := route
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				r = r.WithContext(s.ctx)
				route.Func.ServeHTTP(w, r)
			})
	}
	s.server = httptest.NewServer(router)
}

func (s *ApiSuite) TearDownTest() {
	s.mockCtl.Finish()
}

func (s *ApiSuite) TestAddBannerToSlotSuccess() {
	form := serverapi.BannerSlotForm{BannerID: 1, SlotID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().AddBannerToSlot(s.ctx, form.BannerID, form.SlotID).Return(nil)
	resp, err := http.Post(s.server.URL+"/slot/banner/add", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *ApiSuite) TestAddBannerToSlotInvalidInput() {
	data := []byte("invalid")
	readers := []io.Reader{
		bytes.NewReader(data),
		nil,
	}

	for i, reader := range readers {
		resp, err := http.Post(s.server.URL+"/slot/banner/add", "application/json", reader)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode, "at", i)
	}
}

func (s *ApiSuite) TestAddBannerToSlotBannerNotFound() {
	form := serverapi.BannerSlotForm{BannerID: 1, SlotID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().AddBannerToSlot(s.ctx, form.BannerID, form.SlotID).Return(storage.ErrObjectNotFound)
	resp, err := http.Post(s.server.URL+"/slot/banner/add", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ApiSuite) TestAddBannerToSlotSlotNotFound() {
	form := serverapi.BannerSlotForm{BannerID: 1, SlotID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().AddBannerToSlot(s.ctx, form.BannerID, form.SlotID).Return(storage.ErrObjectNotFound)
	resp, err := http.Post(s.server.URL+"/slot/banner/add", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ApiSuite) TestAddBannerToSlotUnknownError() {
	form := serverapi.BannerSlotForm{BannerID: 1, SlotID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().AddBannerToSlot(s.ctx, form.BannerID, form.SlotID).Return(errUnknown)
	resp, err := http.Post(s.server.URL+"/slot/banner/add", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *ApiSuite) TestRemoveBannerFromSlotSuccess() {
	form := serverapi.BannerSlotForm{BannerID: 1, SlotID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().RemoveBannerFromSlot(s.ctx, form.BannerID, form.SlotID).Return(nil)
	resp, err := http.Post(s.server.URL+"/slot/banner/remove", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *ApiSuite) TestRemoveBannerFromSlotInvalidInput() {
	data := []byte("invalid")
	readers := []io.Reader{
		bytes.NewReader(data),
		nil,
	}

	for i, reader := range readers {
		resp, err := http.Post(s.server.URL+"/slot/banner/remove", "application/json", reader)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode, "at", i)
	}
}

func (s *ApiSuite) TestRemoveBannerFromSlotSlotNotFound() {
	form := serverapi.BannerSlotForm{BannerID: 1, SlotID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().RemoveBannerFromSlot(s.ctx, form.BannerID, form.SlotID).Return(storage.ErrObjectNotFound)
	resp, err := http.Post(s.server.URL+"/slot/banner/remove", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ApiSuite) TestRemoveBannerFromSlotBannerInSlotNotFound() {
	form := serverapi.BannerSlotForm{BannerID: 1, SlotID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().RemoveBannerFromSlot(s.ctx, form.BannerID, form.SlotID).Return(storage.ErrObjectNotFound)
	resp, err := http.Post(s.server.URL+"/slot/banner/remove", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ApiSuite) TestRemoveBannerFromSlotUnknownError() {
	form := serverapi.BannerSlotForm{BannerID: 1, SlotID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().RemoveBannerFromSlot(s.ctx, form.BannerID, form.SlotID).Return(errUnknown)
	resp, err := http.Post(s.server.URL+"/slot/banner/remove", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *ApiSuite) TestBannerForSlotSuccess() {
	query := serverapi.BannerForSlotForm{
		SlotID:   1,
		SocDemID: 1,
	}

	s.mockStore.EXPECT().BannersStatistics(s.ctx, query.SlotID, query.SocDemID).Return(mockStatistics(), nil)
	s.mockStore.EXPECT().AddViewForBanner(s.ctx, mockStatistics()[2].BannerID, query.SlotID, query.SocDemID).Return(nil)
	resp, err := http.Get(s.server.URL + fmt.Sprintf("/banner?slot_id=%d&soc_dem_id=%d", query.SlotID, query.SocDemID))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *ApiSuite) TestBannerForSlotInvalidInput() {
	resp, err := http.Get(s.server.URL + fmt.Sprintf("/banner?id=%d&soc=w", 1))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *ApiSuite) TestBannerForSlotSlotNotFound() {
	query := serverapi.BannerForSlotForm{
		SlotID:   1,
		SocDemID: 1,
	}

	s.mockStore.EXPECT().BannersStatistics(s.ctx, query.SlotID, query.SocDemID).Return(nil, storage.ErrObjectNotFound)
	resp, err := http.Get(s.server.URL + fmt.Sprintf("/banner?slot_id=%d&soc_dem_id=%d", query.SlotID, query.SocDemID))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ApiSuite) TestBannerForSlotSocialGroupNotFound() {
	query := serverapi.BannerForSlotForm{
		SlotID:   1,
		SocDemID: 1,
	}

	s.mockStore.EXPECT().BannersStatistics(s.ctx, query.SlotID, query.SocDemID).Return(nil, storage.ErrObjectNotFound)
	resp, err := http.Get(s.server.URL + fmt.Sprintf("/banner?slot_id=%d&soc_dem_id=%d", query.SlotID, query.SocDemID))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ApiSuite) TestBannerForSlotAddViewError() {
	query := serverapi.BannerForSlotForm{
		SlotID:   1,
		SocDemID: 1,
	}

	s.mockStore.EXPECT().BannersStatistics(s.ctx, query.SlotID, query.SocDemID).Return(mockStatistics(), nil)
	s.mockStore.EXPECT().AddViewForBanner(s.ctx, mockStatistics()[2].BannerID, query.SlotID, query.SocDemID).Return(errors.New("storeErr"))
	resp, err := http.Get(s.server.URL + fmt.Sprintf("/banner?slot_id=%d&soc_dem_id=%d", query.SlotID, query.SocDemID))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *ApiSuite) TestBannerForSlotUnknownError() {
	query := serverapi.BannerForSlotForm{
		SlotID:   1,
		SocDemID: 1,
	}

	s.mockStore.EXPECT().BannersStatistics(s.ctx, query.SlotID, query.SocDemID).Return(nil, errUnknown)
	resp, err := http.Get(s.server.URL + fmt.Sprintf("/banner?slot_id=%d&soc_dem_id=%d", query.SlotID, query.SocDemID))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *ApiSuite) TestAddCLickForBannerSuccess() {
	form := serverapi.BannerClickFrom{BannerID: 1, SlotID: 1, SocDemID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().AddClickForBanner(s.ctx, form.BannerID, form.SlotID, form.SocDemID).Return(nil)
	resp, err := http.Post(s.server.URL+"/banner/click/add", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *ApiSuite) TestAddCLickForBannerInvalidInput() {
	data := []byte("invalid")
	readers := []io.Reader{
		bytes.NewReader(data),
		nil,
	}

	for i, reader := range readers {
		resp, err := http.Post(s.server.URL+"/banner/click/add", "application/json", reader)

		s.Require().NoError(err)
		s.Require().Equal(http.StatusBadRequest, resp.StatusCode, "at", i)
	}
}

func (s *ApiSuite) TestAddCLickForBannerSlotNotFound() {
	form := serverapi.BannerClickFrom{BannerID: 1, SlotID: 1, SocDemID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().AddClickForBanner(s.ctx, form.BannerID, form.SlotID, form.SocDemID).Return(storage.ErrObjectNotFound)
	resp, err := http.Post(s.server.URL+"/banner/click/add", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ApiSuite) TestAddCLickForBannerBannerInSlotNotFound() {
	form := serverapi.BannerClickFrom{BannerID: 1, SlotID: 1, SocDemID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().AddClickForBanner(s.ctx, form.BannerID, form.SlotID, form.SocDemID).Return(storage.ErrObjectNotFound)
	resp, err := http.Post(s.server.URL+"/banner/click/add", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ApiSuite) TestAddCLickForBannerSocialGroupNotFound() {
	form := serverapi.BannerClickFrom{BannerID: 1, SlotID: 1, SocDemID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().AddClickForBanner(s.ctx, form.BannerID, form.SlotID, form.SocDemID).Return(storage.ErrObjectNotFound)
	resp, err := http.Post(s.server.URL+"/banner/click/add", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *ApiSuite) TestAddCLickForBannerUnknownError() {
	form := serverapi.BannerClickFrom{BannerID: 1, SlotID: 1, SocDemID: 1}
	data, err := json.Marshal(&form)

	s.Require().NoError(err)

	s.mockStore.EXPECT().AddClickForBanner(s.ctx, form.BannerID, form.SlotID, form.SocDemID).Return(errUnknown)
	resp, err := http.Post(s.server.URL+"/banner/click/add", "application/json", bytes.NewReader(data))

	s.Require().NoError(err)
	s.Require().Equal(http.StatusBadRequest, resp.StatusCode)
}

func TestApiSuite(t *testing.T) {
	suite.Run(t, new(ApiSuite))
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
