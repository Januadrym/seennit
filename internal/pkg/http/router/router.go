package router

import "net/http"

import "github.com/gorilla/mux"

type (
	// Config hold configurations of router
	Config struct {
		NotFoundHandler http.Handler
		Routes          []Route
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
		r.Path(rt.Path).Methods(rt.Method).Handler(h).Queries(rt.Queries...)
	}
	if conf.NotFoundHandler != nil {
		r.NotFoundHandler = conf.NotFoundHandler
	}

	return r, nil
}
