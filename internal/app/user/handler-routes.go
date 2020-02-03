package user

import (
	"net/http"
	"github.com/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/user/getall",
			Method:  http.MethodGet,
			Handler: h.GetAll,
		},
		{
			Path:    "/user/register",
			Method:  http.MethodPost,
			Handler: h.Register,
		},
		{
			Path:    "/user/getsingle",
			Method:  http.MethodGet,
			Handler: h.GetUsers,
		},
		{
			// just for testing
			Path:    "/user/deletealluser",
			Method:  http.MethodDelete,
			Handler: h.DeleteAll,
		},
	}
}
