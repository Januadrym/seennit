package comment

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/http/respond"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type (
	service interface {
		GetAll(ctx context.Context) ([]*types.Comment, error)
		Update(ctx context.Context, id string, c *types.Comment) error
		DeleteByID(ctx context.Context, id string) error
	}

	Handler struct {
		Svc service
	}
)

func NewHander(svc service) *Handler {
	return &Handler{
		Svc: svc,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.Svc.GetAll(r.Context())
	if err != nil {
		logrus.WithContext(r.Context()).Infof("failed to get all comment, err: %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: list,
	})
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Info("invalid id")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	var cm *types.Comment
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&cm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.Svc.Update(r.Context(), id, cm); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: cm,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Info("invalid id")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	if err := h.Svc.DeleteByID(r.Context(), id); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: id,
		},
	})
}
