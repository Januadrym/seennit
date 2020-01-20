package types

type (
	Response struct {
		Code  string
		Data  interface{}
		Error string
	}
)

var (
	CodeSuccess = "0000"
)
