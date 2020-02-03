package router

import (
	"net/http"

	"github.com/gorilla/mux"
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

func New(conf *Config) (http.handler, error) {
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

// func LoadConfigFromEnv() *Config {
// 	var conf Config
// 	if err := envconfig.Process("", conf); err != nil {
// 		logrus.Errorf("config: unable to load config for %T: %s", conf, err)
// 	}
// 	return &conf
// }
