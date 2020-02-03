package main

import (
	"github.com/sirupsen/logrus"
	"github.com/Januadrym/seennit/internal/app/api"
	"github.com/Januadrym/seennit/internal/pkg/http/server"
)

func main() {
	router, err := api.NewRouter()
	if err != nil {
		logrus.Panic("Cannot initiate router! err: ", err)
	}

	serverConf := server.LoadConfigFromEnv()
	server.ListenAndServe(serverConf, router)
}
