package user

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:        "/users",
			Method:      http.MethodGet,
			Handler:     h.GetAll,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:    "/users/registeration",
			Method:  http.MethodPost,
			Handler: h.Register,
		},
		{
			// just for testing
			Path:        "/users/one",
			Method:      http.MethodGet,
			Handler:     h.GetUsers,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// just for testing
			Path:        "/users",
			Method:      http.MethodDelete,
			Handler:     h.DeleteAll,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
	}
}
