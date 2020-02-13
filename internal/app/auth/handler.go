package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/http/respond"

	"github.com/sirupsen/logrus"
)

type (
	service interface {
		Auth(ctx context.Context, email, password string) (string, *types.User, error)
	}
	Handler struct {
		srv service
	}
)

func NewHandler(srv service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Email    string
		Password string
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Errorf("server error, err: %v", http.StatusInternalServerError)
		respond.Error(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	token, user, err := h.srv.Auth(r.Context(), req.Email, req.Password)
	if err != nil {
		logrus.Errorf("unauthorized, err: %v", http.StatusUnauthorized)
		respond.Error(w, err, http.StatusUnauthorized)
		return
	}
	respond.JSON(w, http.StatusOK, types.BaseResponse{
		Data: map[string]interface{}{
			"token":     token,
			"user_info": user,
		},
	})
}
