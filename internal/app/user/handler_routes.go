package user

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/users",
			Method:  http.MethodGet,
			Handler: h.GetAll,
		},
		{
			Path:    "/users/registeration",
			Method:  http.MethodPost,
			Handler: h.Register,
		},
		{
			Path:        "/users/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.Delete,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/users",
			Method:      http.MethodPut,
			Handler:     h.Update,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// development fuc :)
			Path:    "/users/one",
			Method:  http.MethodGet,
			Handler: h.GetUsers,
			// Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// development fuc :)
			Path:    "/users/remove/all",
			Method:  http.MethodDelete,
			Handler: h.DeleteAllUser,
		},
	}
}
