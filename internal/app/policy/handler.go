package policy

type (
	service interface {
	}

	Handler struct {
		svc service
	}
)

func NewHandler(svc service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// func (h *Handler) Removetest(w http.ResponseWriter, r *http.Request) {
// 	if err := h.svc.RemoveTest(r.Context()); err != nil {
// 		respond.Error(w, err, http.StatusInternalServerError)
// 	}
// 	respond.JSON(w, http.StatusOK, types.BaseResponse{
// 		Data: "c",
// 	})
// }
