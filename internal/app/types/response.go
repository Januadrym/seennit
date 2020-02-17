package types

import (
	"encoding/json"

	"github.com/Januadrym/seennit/internal/app/status"
)

type (
	BaseResponse struct {
		status.Status
		Data interface{} `json:"data"`
	}
	baseResponse BaseResponse

	IDResponse struct {
		ID string `json:"id"`
	}
)

func (rs BaseResponse) MarshalJSON() ([]byte, error) {
	var v = baseResponse(rs)
	if v.Status.Status() == 0 {
		v.Status = status.Gen().Success
	}
	return json.Marshal(v)
}
