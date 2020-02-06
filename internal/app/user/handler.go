package user

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
		Register(ctx context.Context, req *types.RegisterRequest) (*types.User, error)
		SearchUser(ctx context.Context, req *types.User) (*types.User, error)
		FindAll(ctx context.Context) ([]*types.User, error)
		DeleteAll(ctx context.Context) error
		Delete(ctx context.Context, userID string) error
		Update(ctx context.Context, userID string, user *types.User) error
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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req types.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	usr, err := h.Svc.Register(r.Context(), &req)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}

	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: usr.Strip(),
	})
}

//testing func: to find user
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req *types.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	usr, err := h.Svc.SearchUser(r.Context(), req)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: usr.Strip(),
	})
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.Svc.FindAll(r.Context())
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: users,
	})
}

func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := h.Svc.DeleteAll(r.Context())
	if err != nil {
		respond.JSON(w, http.StatusInternalServerError, types.BaseResponse{
			Error: "server error",
		})
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: "OK",
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Info("invalid id")
		respond.Error(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := h.Svc.Delete(r.Context(), id); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: id,
		},
	})
}
