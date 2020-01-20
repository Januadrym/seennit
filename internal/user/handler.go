package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	service interface {
		Create(context.Context, User) (string, error)
		Home(context.Context, User) (string, error)
		Register(ctx context.Context, req RegisterRequest) (string, error)
	}

	Handler struct {
		Svc service
	}
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	id, err := h.Svc.Create(context.Background(), User{})
	if err != nil {
		fmt.Println("Handle error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something error in server"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	uname, err := h.Svc.Home(context.Background(), User{})
	if err != nil {
		fmt.Println("Handle error", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something error in server"))
		return
	}
	fmt.Println(r.URL)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(uname))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userName, err := h.Svc.Register(r.Context(), req)
	if err == ErrUserAlreadyExist {
		fmt.Println("user already exist")
		return
	}
	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(userName)
	w.Write([]byte(userName))
}
