package post

import (
	"context"
	"encoding/json"
	"fmt"
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
		GetEntire(ctx context.Context) ([]*types.Post, error)
		Create(ctx context.Context, req *types.Post, nameComm string) (*types.Post, error)
		FindByID(ctx context.Context, id string) (*types.Post, error)
		GetAll(ctx context.Context, nameComm string) ([]*types.Post, error)
		UpdatePost(ctx context.Context, id string, p *types.PostUpdateRequest) error
		ChangeStatus(ctx context.Context, id string, status types.Status) error
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
	var req types.Post
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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Infof("invalid id")
		respond.Error(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var req types.PostUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	if err := h.Svc.UpdatePost(r.Context(), id, &req); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: id,
		},
	})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Infof("invalid id")
		respond.Error(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}

	p, err := h.Svc.FindByID(r.Context(), id)
	if err != nil {
		logrus.WithContext(r.Context()).Errorf("post cannot be found, err: %v", err)
		respond.Error(w, err, http.StatusBadRequest)
		return
	}
	if p.Status != types.StatusDelete {
		respond.JSON(w, http.StatusOK, types.BaseResponse{
			Data: p,
		})
		return
	}
	// display draft post for the owner
	user := auth.FromContext(r.Context())
	if p.Status == types.StatusDraft && user != nil && user.UserID == p.CreatedByID {
		respond.JSON(w, http.StatusOK, types.BaseResponse{
			Data: p,
		})
		return
	}
	respond.JSON(w, http.StatusNotFound, types.BaseResponse{
		Status: status.Gen().NotFound,
		Data:   "post not available",
	})
	return
}

//Archive Post: post can no longer be edited or commented
func (h *Handler) ArchivePost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Infof("invalid id")
		respond.Error(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := h.Svc.ChangeStatus(r.Context(), id, types.StatusArchived); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: id,
		},
	})
}

//for futher delete all posts have status: deleted
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Infof("invalid id")
		respond.Error(w, fmt.Errorf("invalid id"), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := h.Svc.ChangeStatus(r.Context(), id, types.StatusDelete); err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: types.IDResponse{
			ID: id,
		},
	})
}

func (h *Handler) GetEntireThing(w http.ResponseWriter, r *http.Request) {
	list, err := h.Svc.GetEntire(r.Context())
	if err != nil {
		logrus.WithContext(r.Context()).Errorf("failed to find, err: %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: list,
	})
}
