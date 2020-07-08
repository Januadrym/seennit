package main

import (
	"github.com/Januadrym/seennit/internal/app/api"
	"github.com/Januadrym/seennit/internal/pkg/http/server"

	"github.com/sirupsen/logrus"
)

func main() {
	router, err := api.NewRouter()
	if err != nil {
		logrus.Panic("Cannot initiate router! err: ", err)
	}

	serverConf := server.LoadConfigFromEnv()
	server.ListenAndServe(serverConf, router)
}
