package server

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	//HTTP server configurations
	Config struct {
		Address      string        `envconfig:"HTTP_ADDRESS" default:":"`
		Port         int           `envconfig:"HTTP_PORT" default:"8080"`
		ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"10s"`
		WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"10s"`
	}
)

func Load() (Config, error) {
	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		return conf, err
	}
	return conf, nil
}
