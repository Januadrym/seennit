package main

import (
	"github.com/sirupsen/logrus"
	"vnmquan.com/seennit/internal/app/api"
	"vnmquan.com/seennit/internal/pkg/http/server"
)

func main() {
	router, err := api.NewRouter()
	if err != nil {
		logrus.Panic("Cannot initiate router! err: ", err)
	}

	serverConf := server.LoadConfigFromEnv()
	server.ListenAndServe(serverConf, router)
}
