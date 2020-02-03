package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/Januadrym/seennit/internal/pkg/config/env"
)

type (
	Config struct {
		Routes          []Route
		NotFoundHandler http.Handler
	}

	Route struct {
		Desc    string
		Path    string
		Method  string
		Queries []string
		Handler http.HandlerFunc
	}
)

func New(conf *Config) (http.Handler, error) {
	r := mux.NewRouter()
	for _, rt := range conf.Routes {
		var h http.Handler
		h = http.HandlerFunc(rt.Handler)

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
