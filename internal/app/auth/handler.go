package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Januadrym/seennit/internal/app/types"
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
		logrus.Errorf("server error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	token, user, err := h.srv.Auth(r.Context(), req.Email, req.Password)
	if err != nil {
		logrus.Errorf("unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(types.Response{
		Data: map[string]interface{}{
			"token: ":     token,
			"user Info: ": user,
		},
	})

}
