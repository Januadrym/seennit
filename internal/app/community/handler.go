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
		SearchCommunity(ctx context.Context, name string) (*Community, error)
		CreateCommunity(ctx context.Context, req *Community) (*Community, error)
		DeleteCommunity(ctx context.Context, idCom string) error
		GetCommunity(ctx context.Context, name string) (*Community, error)
		EnrollUser(ctx context.Context, idCom string) error
		GetAll(ctx context.Context) ([]*Community, error)
		UpdateInfo(ctx context.Context, idCom string, comm *Community) error
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

func (h *Handler) EnrollUser(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	if name == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, fmt.Errorf("invalid name"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	com, err := h.Svc.SearchCommunity(r.Context(), name)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	if err := h.Svc.EnrollUser(r.Context(), com.ID); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: com,
	})
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.Svc.GetAll(r.Context())
	if err != nil {
		logrus.WithContext(r.Context()).Errorf("failed to get all communities, err: %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: list,
	})
}

func (h *Handler) UpdateInfo(w http.ResponseWriter, r *http.Request) {
	var req *Community
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comName := mux.Vars(r)["name"]
	if comName == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, fmt.Errorf("invalid name"), http.StatusBadRequest)
		return
	}
	com, err := h.Svc.SearchCommunity(r.Context(), comName)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	if err := h.Svc.UpdateInfo(r.Context(), com.ID, req); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: com,
	})
}
