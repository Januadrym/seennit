package post

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
		Create(ctx context.Context, req *Post, nameComm string) (*Post, error)
		GetAll(ctx context.Context, nameComm string) ([]*Post, error)
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

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req Post
	comName := mux.Vars(r)["name"]
	if comName == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, fmt.Errorf("invalid name"), http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	post, err := h.Svc.Create(r.Context(), &req, comName)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: post,
	})
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	comName := mux.Vars(r)["name"]
	if comName == "" {
		logrus.WithContext(r.Context()).Info("invalid name")
		respond.Error(w, fmt.Errorf("invalid name"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	posts, err := h.Svc.GetAll(r.Context(), comName)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: posts,
	})
}
