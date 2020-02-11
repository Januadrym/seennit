package community

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:        "/communities",
			Method:      http.MethodPost,
			Handler:     h.CreateCommunity,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/communities",
			Method:      http.MethodGet,
			Handler:     h.SearchTest,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/communities/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.DeleteComByID,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/communities/{name:[a-z0-9-\\-]+}",
			Method:      http.MethodGet,
			Handler:     h.GetCommunity,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
	}
}
