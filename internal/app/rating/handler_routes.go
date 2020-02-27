package rating

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:    "/api/rating",
			Method:  http.MethodPost,
			Handler: h.Create,
		},
	}
}
