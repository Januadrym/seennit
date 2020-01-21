package main

import (
	"fmt"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	mongodbconf "vnmquan.com/seennit/internal/app/config/mongodb"
	serverconf "vnmquan.com/seennit/internal/app/config/server"
	"vnmquan.com/seennit/internal/app/user"
)

func main() {
	conf, err := serverconf.Load()
	if err != nil {
		panic(err)
	}
	router := mux.NewRouter()

	dbconf, err := mongodbconf.Load()
	if err != nil {
		panic(err)
	}

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    dbconf.Addrs,
		Database: dbconf.Database,
		Username: dbconf.Username,
		Password: dbconf.Password,
	})
	if err != nil {
		panic(err)
	}

	userRepo := user.NewMongoDBRepo(session)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	router.Path("/register").Methods(http.MethodPost).HandlerFunc(userHandler.Register)
	router.Path("/").Methods(http.MethodGet).HandlerFunc(userHandler.Home)
	router.Path("/alluser").Methods(http.MethodGet).HandlerFunc(userHandler.GetUsers)

	server := http.Server{
		Addr:         conf.Addr,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		Handler:      router,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}

}
