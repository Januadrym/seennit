package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/http/respond"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type (
	service interface {
		Register(ctx context.Context, req *types.RegisterRequest) (*types.User, error)
		FindAll(ctx context.Context) ([]*types.User, error)
		Delete(ctx context.Context, userID string) error
		Update(ctx context.Context, userID string, user *types.User) error

		LoadNotification(ctx context.Context, ID string) ([]*types.PushNotification, error)
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

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Info("invalid id")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var user *types.User
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authen := auth.FromContext(r.Context())
	err := h.Svc.Update(r.Context(), authen.UserID, user)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: user,
	})
}

func (h *Handler) LoadNoti(w http.ResponseWriter, r *http.Request) {
	thisUser := auth.FromContext(r.Context())

	notis, err := h.Svc.LoadNotification(r.Context(), thisUser.UserID)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: notis,
	})
}
