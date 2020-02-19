package community

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		// s stand for sub, which mean a community
		{
			// create community
			Path:        "/s",
			Method:      http.MethodPost,
			Handler:     h.CreateCommunity,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// get all community
			Path:        "/s",
			Method:      http.MethodGet,
			Handler:     h.GetAll,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}",
			Method:      http.MethodGet,
			Handler:     h.GetCommunity,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// temporary func for delete a community
			// should have status: active, hidden, archive
			Path:   "/s/{name:[a-z0-9-A-Z-\\-]+}",
			Method: http.MethodDelete,
			// Handler:     h.DeleteComByID,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},

		{
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/enrollment",
			Method:      http.MethodPut,
			Handler:     h.EnrollUser,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}",
			Method:      http.MethodPut,
			Handler:     h.UpdateInfo,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		// change status to hidden
		{
			Path:   "/s/{name:[a-z0-9-A-Z-\\-]+}",
			Method: http.MethodPatch,
			// Handler:     h.DeleteCommunity,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
	}
}
