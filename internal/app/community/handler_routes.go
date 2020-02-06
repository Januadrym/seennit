package community

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/communities/create",
			Method:  http.MethodPost,
			Handler: h.CreateCommunity,
		},
	}
}
