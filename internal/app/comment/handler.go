package comment

import (
	"context"
	"net/http"

	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/http/respond"
	"github.com/sirupsen/logrus"
)

type (
	service interface {
		GetAll(ctx context.Context) ([]*types.Comment, error)
	}

	Handler struct {
		Svc service
	}
)

func NewHander(svc service) *Handler {
	return &Handler{
		Svc: svc,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.Svc.GetAll(r.Context())
	if err != nil {
		logrus.WithContext(r.Context()).Infof("failed to get all comment, err: %v", err)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: list,
	})
}

// func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {

// }
