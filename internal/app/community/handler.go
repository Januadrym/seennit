package community

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/http/respond"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type (
	service interface {
		SearchCommunity(ctx context.Context, req *Community) (*Community, error)
		CreateCommunity(ctx context.Context, req *Community) (*Community, error)
		DeleteCommunity(ctx context.Context, comID string) error
		GetCommunity(ctx context.Context, name string) (*Community, error)
	}

	Handler struct {
		Svc service
	}
)

func NewHandler(svc service) *Handler {
	return &Handler{
		Svc: svc,
	}
}

func (h *Handler) CreateCommunity(w http.ResponseWriter, r *http.Request) {
	var req Community
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	comm, err := h.Svc.CreateCommunity(r.Context(), &req)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: comm,
	})
}

func (h *Handler) SearchTest(w http.ResponseWriter, r *http.Request) {
	var req Community
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	com, err := h.Svc.SearchCommunity(r.Context(), &req)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: com,
	})
}

func (h *Handler) DeleteComByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Info("invalid id")
		respond.Error(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := h.Svc.DeleteCommunity(r.Context(), id); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: id,
		},
	})
}

func (h *Handler) GetCommunity(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, fmt.Errorf("invalid name"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	com, err := h.Svc.GetCommunity(r.Context(), name)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: com,
	})
}
