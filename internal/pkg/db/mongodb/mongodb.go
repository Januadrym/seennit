package mongodb

import (
	"time"

	"github.com/globalsign/mgo"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
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

func Load() (*Config, error) {
	var conf Config
	if err := envconfig.Process("", &conf); err != nil {
		return &conf, err
	}
	return &conf, nil
}

//// envconfig.go
//// for loading config from env
// func LoadConfigFromEnv() *Config {
// 	var conf Config
// 	envconfig.Load(&conf)
// 	return &conf
// }

func Dial(conf *Config) (*mgo.Session, error) {
	logrus.Infof("dialing to target MongoDB at: %v, database: %v", conf.Addrs, conf.Database)
	ms, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    conf.Addrs,
		Database: conf.Database,
		Username: conf.Username,
		Password: conf.Password,
		Timeout:  conf.Timeout,
	})
	if err != nil {
		return nil, err
	}
	return ms, nil
}

func (conf *Config) DialInfo() *mgo.DialInfo {
	return &mgo.DialInfo{
		Addrs:    conf.Addrs,
		Database: conf.Database,
		Username: conf.Username,
		Password: conf.Password,
		Timeout:  conf.Timeout,
	}
}
