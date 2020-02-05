package respond

import (
	"encoding/json"
	"net/http"
)

type (
	appError interface {
		Error() string
		Code() uint32
		Message() string
		Status() int
	}
)

func JSON(w http.ResponseWriter, status int, data interface{}) {
	b, err := json.Marshal(data) // byte
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}

func Error(w http.ResponseWriter, err error, status int) {
	if appError, ok := err.(appError); ok {
		JSON(w, appError.Status(), appError)
		return
	}
	JSON(w, status, map[string]interface{}{
		"code":    status,
		"message": err.Error(),
	})
}
