package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/nsmak/bannersRotation/internal/app"
	"github.com/nsmak/bannersRotation/internal/server/rest"
	"github.com/nsmak/bannersRotation/internal/storage"
)

type BannerSlotForm struct {
	BannerID int64 `json:"banner_id"`
	SlotID   int64 `json:"slot_id"`
}

type BannerForSlotForm struct {
	SlotID   int64 `schema:"slot_id"`
	SocDemID int64 `schema:"soc_dem_id"`
}

type BannerClickFrom struct {
	BannerID int64 `json:"banner_id"`
	SlotID   int64 `json:"slot_id"`
	SocDemID int64 `json:"soc_dem_id"`
}

type API struct {
	rotator *app.RotatorDomain
}

func New(rotator *app.RotatorDomain) *API {
	return &API{rotator: rotator}
}

func (a *API) addBannerToSlot(w http.ResponseWriter, r *http.Request) {
	var form BannerSlotForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse")
		return
	}

	if err := a.rotator.AddBannerToSlot(r.Context(), form.BannerID, form.SlotID); err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, storage.ErrObjectNotFound) {
			statusCode = http.StatusNotFound
		}
		rest.SendErrorJSON(w, r, statusCode, err, "")
		return
	}

	rest.SendDataJSON(w, r, http.StatusOK, nil)
}

func (a *API) removeBannerFromSlot(w http.ResponseWriter, r *http.Request) {
	var form BannerSlotForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse")
		return
	}

	if err := a.rotator.RemoveBannerFromSlot(r.Context(), form.BannerID, form.SlotID); err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, storage.ErrObjectNotFound) {
			statusCode = http.StatusNotFound
		}
		rest.SendErrorJSON(w, r, statusCode, err, "")
		return
	}

	rest.SendDataJSON(w, r, http.StatusOK, nil)
}

func (a *API) bannerForSlot(w http.ResponseWriter, r *http.Request) {
	var query BannerForSlotForm
	if err := schema.NewDecoder().Decode(&query, r.URL.Query()); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't get query params")
		return
	}

	bannerID, err := a.rotator.BannerIDForSlot(r.Context(), query.SlotID, query.SocDemID)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, storage.ErrObjectNotFound) {
			statusCode = http.StatusNotFound
		}
		rest.SendErrorJSON(w, r, statusCode, err, "can't get banner id")
		return
	}

	rest.SendDataJSON(w, r, http.StatusOK, rest.JSON{"banner_id": bannerID})
}

func (a *API) addCLickForBanner(w http.ResponseWriter, r *http.Request) {
	var form BannerClickFrom
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse")
		return
	}

	err := a.rotator.AddClickForBanner(r.Context(), form.BannerID, form.SlotID, form.SocDemID)
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, storage.ErrObjectNotFound) {
			statusCode = http.StatusNotFound
		}
		rest.SendErrorJSON(w, r, statusCode, err, "")
		return
	}

	rest.SendDataJSON(w, r, http.StatusOK, nil)
}

func (a *API) Routes() []rest.Route {
	return []rest.Route{
		{
			Name:   "AddBannerToSlot",
			Method: http.MethodPost,
			Path:   "/slot/banner/add",
			Func:   a.addBannerToSlot,
		},
		{
			Name:   "RemoveBannerFromSlot",
			Method: http.MethodPost,
			Path:   "/slot/banner/remove",
			Func:   a.removeBannerFromSlot,
		},
		{
			Name:   "BannerForSlot",
			Method: http.MethodGet,
			Path:   "/banner",
			Func:   a.bannerForSlot,
		},
		{
			Name:   "ClickForBanner",
			Method: http.MethodPost,
			Path:   "/banner/click/add",
			Func:   a.addCLickForBanner,
		},
	}
}
