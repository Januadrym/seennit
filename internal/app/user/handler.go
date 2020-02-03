package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/seennit/internal/app/types"
)

type (
	service interface {
		Register(ctx context.Context, req *types.RegisterRequest) (*types.User, error)
		SearchUser(ctx context.Context, req *types.User) (*types.User, error)
		FindAll(ctx context.Context) ([]*types.User, error)
		DeleteAll(ctx context.Context) error
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
	usr, err := h.Svc.Register(r.Context(), &req)
	if err == ErrUserAlreadyExist {
		json.NewEncoder(w).Encode(types.Response{
			Code:  "0002",
			Error: err.Error(),
		})
		return
	}
	if err != nil {
		json.NewEncoder(w).Encode(types.Response{
			Code: "0001",
		})
		return
	}

	json.NewEncoder(w).Encode(types.Response{
		Code: types.CodeSuccess,
		Data: usr.Strip(),
	})
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req *types.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	urs, err := h.Svc.SearchUser(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something error in server"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.Response{
		Code: types.CodeSuccess,
		Data: urs.Strip(),
	})
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.Svc.FindAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something error in server"))
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.Response{
		Code: types.CodeSuccess,
		Data: users,
	})
}

func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := h.Svc.DeleteAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something error in server"))
		return
	}
	w.WriteHeader(http.StatusOK)
}
