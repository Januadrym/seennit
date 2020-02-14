package post

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/submit",
			Method:      http.MethodPost,
			Handler:     h.CreatePost,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// get all post of one community
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}",
			Method:      http.MethodGet,
			Handler:     h.GetAll,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
	}
}
