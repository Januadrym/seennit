package mongodbconf

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		Addrs    []string      `envconfig:"MONGODB_ADDRS" default:"127.0.0.1:27017"`
		Database string        `envconfig:"MONGODB_DATABASE" default:"notes"`
		Username string        `envconfig:"MONGODB_USERNAME"`
		Password string        `envconfig:"MONGODB_PASSWORD"`
		Timeout  time.Duration `envconfig:"MONGODB_TIMEOUT" default:"10s"`
	}
)

func Load() (Config, error) {
	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		return conf, err
	}
	return conf, nil
}
