package comment

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			// get all comments of all post (all comment that ever exist)
			Path:    "/comments",
			Method:  http.MethodGet,
			Handler: h.GetAll,
		},
		{
			Path:        "/comments/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodPut,
			Handler:     h.Update,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/comments/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.Delete,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
	}
}
