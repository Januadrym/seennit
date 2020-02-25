package policy

import (
	"context"
	"net/http"

	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/http/respond"
)

type (
	service interface {
		GetAllMods(ctx context.Context) error
	}

	Handler struct {
		svc service
	}
)

func NewHandler(svc service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) GetMods(w http.ResponseWriter, r *http.Request) {
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: "hello ok",
	})
	err := h.svc.GetAllMods(r.Context())
	if err != nil {
		return
	}
	// respond.JSON(w, http.StatusOK, types.BaseResponse{
	// 	Data: "hello ok",
	// })
}
