package post

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			// get one post
			Path:        "/posts/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodGet,
			Handler:     h.GetPost,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/posts/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodPut,
			Handler:     h.UpdatePost,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// change status to archive
			Path:        "/posts/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodPatch,
			Handler:     h.ArchivePost,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// change status to delete
			Path:        "/posts/{id:[a-z0-9-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.DeletePost,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},

		// Comments
		{
			Path:        "/posts/{id:[a-z0-9-\\-]+}/comments",
			Method:      http.MethodPost,
			Handler:     h.CreateComment,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// get all comments of 1 post
			Path:        "/posts/{id:[a-z0-9-\\-]+}/comments",
			Method:      http.MethodGet,
			Handler:     h.GetCommentsPost,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
	}
}
