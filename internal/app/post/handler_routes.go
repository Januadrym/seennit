package post

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		// s stand for sub, which mean a community
		{
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/posts",
			Method:      http.MethodPost,
			Handler:     h.CreatePost,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// get all post of a community
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/posts",
			Method:      http.MethodGet,
			Handler:     h.GetAll,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodPut,
			Handler:     h.Update,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:    "/s/{name:[a-z0-9-A-Z-\\-]+}/{id:[a-z0-9-\\-]+}",
			Method:  http.MethodGet,
			Handler: h.Get,
		},
		{
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodPatch,
			Handler:     h.ArchivePost, // change status to archive
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.Delete, // change status to delete
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// get all post in all community
			Path:    "/posts",
			Method:  http.MethodGet,
			Handler: h.GetEntireThing,
		},
	}
}
