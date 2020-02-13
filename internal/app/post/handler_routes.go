package post

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:        "/posts",
			Method:      http.MethodPost,
			Handler:     h.CreatePost,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
	}
}
