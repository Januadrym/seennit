package router

import (
	"net/http"

	"github.com/Januadrym/seennit/internal/pkg/config/env"
	"github.com/gorilla/mux"
)

type (
	Config struct {
		Middlewares     []Middleware
		Routes          []Route
		NotFoundHandler http.Handler
	}

	Middleware = func(http.Handler) http.Handler

	Route struct {
		Desc        string
		Path        string
		Method      string
		Queries     []string
		Handler     http.HandlerFunc
		Middlewares []Middleware
	}
)

func New(conf *Config) (http.Handler, error) {
	r := mux.NewRouter()
	for _, middleware := range conf.Middlewares {
		r.Use(middleware)
	}

	for _, rt := range conf.Routes {
		var h http.Handler
		h = http.HandlerFunc(rt.Handler)
		for i := len(rt.Middlewares) - 1; i >= 0; i-- {
			h = rt.Middlewares[i](h)
		}
		r.Path(rt.Path).Methods(rt.Method).Handler(h).Queries(rt.Queries...)
	}

	if conf.NotFoundHandler != nil {
		r.NotFoundHandler = conf.NotFoundHandler
	}

	return r, nil
}

func LoadConfigFromEnv() *Config {
	var conf Config
	env.Load(&conf)
	return &conf
}
