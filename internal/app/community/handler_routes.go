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
			Handler:     h.GetAllCommunity,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}",
			Method:      http.MethodGet,
			Handler:     h.GetCommunity,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// change status to private
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}",
			Method:      http.MethodDelete,
			Handler:     h.DeleteCom,
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
		// Policy stuff
		{
			//promote user to moderator
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/users",
			Method:      http.MethodPost,
			Handler:     h.PromoteMod,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// get all enrolled users
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/users",
			Method:      http.MethodGet,
			Handler:     h.GetUsers,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// get all moderators
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/policy",
			Method:      http.MethodGet,
			Handler:     h.GetAllMods,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// Transfer Admin
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/policy",
			Method:      http.MethodPost,
			Handler:     h.TransferAdmin,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},

		// Posts stuff
		{
			// submit one post
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/posts",
			Method:      http.MethodPost,
			Handler:     h.SubmitPost,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
		{
			// get all post of one community
			Path:        "/s/{name:[a-z0-9-A-Z-\\-]+}/posts",
			Method:      http.MethodGet,
			Handler:     h.GetAllPosts,
			Middlewares: []router.Middleware{auth.RequireAuthMiddleware},
		},
	}
}
