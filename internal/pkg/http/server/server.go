package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
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

	// graceful shutdown

	// signs := make(chan os.Signal, 1)
	// signal.Notify(signs, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// <-signs

	if err := srv.ListenAndServe(); err != nil {
		logrus.Panicf("listen: %s\n", err)
	}

	//
	// srvContext, srvCancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
	// defer srvCancel()
	// logrus.Infof("shutting down server in slow motion... ")
	// if err := srv.Shutdown(srvContext); err != nil {
	// 	logrus.Panic("Server is shutting down with error: ", err)
	// }

}
