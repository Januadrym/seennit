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
