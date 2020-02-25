package post

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
		GetEntire(ctx context.Context) ([]*types.Post, error)

		FindByID(ctx context.Context, idPost string) (*types.Post, error)
		UpdatePost(ctx context.Context, idPost string, p *types.PostUpdateRequest) error
		ChangeStatus(ctx context.Context, idPost string, stat types.Status) error

		CreateComment(ctx context.Context, req *types.Comment, idPost string) (*types.Comment, error)
		GetCommentsPost(ctx context.Context, idPost string) ([]*types.Comment, error)
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

// GetEntireThing get all post to display in homepage
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

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Infof("invalid id")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
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
	respond.JSON(w, http.StatusNotFound, types.BaseResponse{
		Status: status.Gen().NotFound,
	})
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Infof("invalid id")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
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

// ArchivePost post can no longer be edited or commented
func (h *Handler) ArchivePost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Infof("invalid id")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
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

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		logrus.WithContext(r.Context()).Infof("invalid id")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
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

// Comment

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req *types.Comment
	idPost := mux.Vars(r)["id"]
	if idPost == "" {
		logrus.WithContext(r.Context()).Info("invalid id")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	comment, err := h.Svc.CreateComment(r.Context(), req, idPost)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: comment,
	})
}

func (h *Handler) GetCommentsPost(w http.ResponseWriter, r *http.Request) {
	idPost := mux.Vars(r)["id"]
	if idPost == "" {
		logrus.WithContext(r.Context()).Info("invalid id")
		respond.Error(w, status.Gen().BadRequest, http.StatusBadRequest)
		return
	}
	comments, err := h.Svc.GetCommentsPost(r.Context(), idPost)
	if err != nil {
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: comments,
	})
}
