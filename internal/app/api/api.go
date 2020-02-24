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
	// Policy
	policySrv, err := newPolicyService()
	if err != nil {
		return nil, err
	}

	// User
	userSrv, err := newUserService()
	if err != nil {
		return nil, err
	}
	userHandler := newUserHandler(userSrv)

	// Comment
	commentSrv, err := newCommentService(policySrv)
	if err != nil {
		return nil, err
	}
	commentHandler := newCommentHandler(commentSrv)

	// Post
	postSrv, err := newPostService(policySrv, commentSrv)
	if err != nil {
		return nil, err
	}
	postHandler := newPostHandler(postSrv)

	// Community
	commSrv, err := newCommunityService(policySrv, postSrv, userSrv)
	if err != nil {
		return nil, err
	}
	commHandler := newCommunityHandler(commSrv)

	jwtSignVerifier := newJWTSignVerifier()
	authHandler := newAuthHandler(jwtSignVerifier, userSrv)
	userInfoMiddleware := auth.UserInfoMiddleware(jwtSignVerifier)

	routes := []router.Route{
		{
			Path:    "/",
			Method:  m_get,
			Handler: postHandler.GetEntireThing,
		},
	}

	routes = append(routes, userHandler.Routes()...)
	routes = append(routes, authHandler.Routes()...)
	routes = append(routes, commHandler.Routes()...)
	routes = append(routes, postHandler.Routes()...)
	routes = append(routes, commentHandler.Routes()...)

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
