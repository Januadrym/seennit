package api

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/pkg/http/router"
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

	conf := router.LoadConfigFromEnv()
	conf.Routes = routes

	r, err := router.New(conf)
	if err != nil {
		return nil, err
	}

	return r, nil
}
