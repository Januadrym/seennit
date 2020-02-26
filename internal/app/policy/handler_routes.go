package policy

import (
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

func (h *Handler) Routes() []router.Route {
	return []router.Route{
		// {
		// 	Path:    "/api/test",
		// 	Method:  http.MethodPost,
		// 	Handler: h.Removetest,
		// },
	}
}
