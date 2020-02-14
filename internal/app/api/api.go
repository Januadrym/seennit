package api

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/pkg/http/middleware"
	"github.com/Januadrym/seennit/internal/pkg/http/router"
)

const (
	m_get     = http.MethodGet
	m_post    = http.MethodPost
	m_put     = http.MethodPut
	m_delete  = http.MethodDelete
	m_options = http.MethodOptions
)

func NewRouter() (http.Handler, error) {
	//Policy
	policySrv, err := newPolicyService()
	if err != nil {
		return nil, err
	}

	//User
	userSrv, err := newUserService()
	if err != nil {
		return nil, err
	}
	userHandler := newUserHandler(userSrv)

	//Community
	commSrv, err := newCommunityService(policySrv)
	if err != nil {
		return nil, err
	}
	commHandler := newCommunityHandler(commSrv)

	//Post
	postSrv, err := newPostService(policySrv, commSrv)
	if err != nil {
		return nil, err
	}
	postHandler := newPostHandler(postSrv)

	jwtSignVerifier := newJWTSignVerifier()
	authHandler := newAuthHandler(jwtSignVerifier, userSrv)
	userInfoMiddleware := auth.UserInfoMiddleware(jwtSignVerifier)

	routes := []router.Route{
		{
			Path:    "/",
			Method:  m_get,
			Handler: ServeHTTP,
		},
	}

	routes = append(routes, userHandler.Routes()...)
	routes = append(routes, authHandler.Routes()...)
	routes = append(routes, commHandler.Routes()...)
	routes = append(routes, postHandler.Routes()...)

	conf := router.LoadConfigFromEnv()
	conf.Routes = routes
	conf.Middlewares = []router.Middleware{
		userInfoMiddleware,
	}

	r, err := router.New(conf)
	if err != nil {
		return nil, err
	}

	return middleware.CORS(r), nil
}
