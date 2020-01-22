package server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func ListenAndServe(conf Config, router http.Handler) {
	port := fmt.Sprint(conf.Port)
	if conf.Port == 0 {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}
	address := fmt.Sprintf("%s:%s", conf.Address, port)
	srv := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	}
	logrus.Infof("HTTP server is listening on: %s", address)

	if err := srv.ListenAndServe(); err != nil {
		logrus.Panicf("listen: %s\n", err)
	}
}
