package api

import (
	"net/http"

	"vnmquan.com/seennit/internal/pkg/http/router"
)

const (
	get     = http.MethodGet
	post    = http.MethodPost
	put     = http.MethodPut
	delete  = http.MethodDelete
	options = http.MethodOptions
)

func NewRouter() (http.Handler, error) {

	//User
	userSrv, err := newUserService()
	if err != nil {
		return nil, err
	}
	userHandler := newUserHandler(userSrv)
	//

	// indexHandler :=
	routes := []router.Route{
		{
			Path:    "/",
			Method:  get,
			Handler: ServeHTTP,
		},
	}

	routes = append(routes, userHandler.Routes()...)

	r, err := router.New()
	if err != nil {
		return nil, err
	}

	return r, nil
}
