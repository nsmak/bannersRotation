package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nsmak/bannersRotation/internal/app"
	"github.com/nsmak/bannersRotation/internal/server/rest"
	"github.com/nsmak/bannersRotation/internal/storage"
)

type BannerSlotForm struct {
	BannerID int64 `json:"banner_id"`
	SlotID   int64 `json:"slot_id"`
}

type API struct {
	rotator *app.RotatorDomain
}

func New(rotator *app.RotatorDomain) *API {
	return &API{rotator: rotator}
}

func (a API) addBannerToSlot(w http.ResponseWriter, r *http.Request) {
	var form BannerSlotForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse")
		return
	}

	if err := a.rotator.AddBannerToSlot(r.Context(), form.BannerID, form.SlotID); err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, storage.ErrBannerNotFound) || errors.Is(err, storage.ErrSlotNotFound) {
			statusCode = http.StatusNotFound
		}
		rest.SendErrorJSON(w, r, statusCode, err, "can't update event")
		return
	}

	rest.SendDataJSON(w, r, http.StatusOK, nil)
}

func (a API) removeBannerFromSlot(w http.ResponseWriter, r *http.Request) {
	var form BannerSlotForm
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		rest.SendErrorJSON(w, r, http.StatusBadRequest, err, "can't parse")
		return
	}

	if err := a.rotator.RemoveBannerFromSlot(r.Context(), form.BannerID, form.SlotID); err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, storage.ErrBannerInSlotNotFound) {
			statusCode = http.StatusNotFound
		}
		rest.SendErrorJSON(w, r, statusCode, err, "can't update event")
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
	}
}
