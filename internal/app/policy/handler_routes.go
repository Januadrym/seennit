package policy

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		{
			Path:   "/api/v1",
			Method: http.MethodGet,
			// Handler: h.GetMods,
		},
	}
}
