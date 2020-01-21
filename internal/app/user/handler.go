package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"vnmquan.com/seennit/internal/app/types"
)

type (
	service interface {
		// Create(context.Context, User) (string, error)
		Home(context.Context, User) (string, error)
		Register(ctx context.Context, req RegisterRequest) (string, error)
		SearchUser(ctx context.Context, req User) (string, error)
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

// func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
// 	id, err := h.Svc.Create(context.Background(), User{})
// 	if err != nil {
// 		fmt.Println("Handle error", err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte("Something error in server"))
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	w.Write([]byte(id))
// }

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
	id, err := h.Svc.Register(r.Context(), req)
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
		Data: map[string]interface{}{
			"id":    id,
			"name":  req.FirstName,
			"email": req.Email,
		},
	})
	// json.NewEncoder(w).Encode(id)
	w.Write([]byte(req.FirstName))
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ursName, err := h.Svc.SearchUser(r.Context(), req)
	if err != nil {
		fmt.Errorf("err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something error in server"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(ursName))
}
