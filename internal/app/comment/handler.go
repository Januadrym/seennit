package comment

type (
	service interface {
	}

	Handler struct {
		Svc service
	}
)

func NewHander(svc service) *Handler {
	return &Handler{
		Svc: svc,
	}
}
