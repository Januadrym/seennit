package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"vnmquan.com/seennit/internal/app/user"
)

func main() {
	userRepo := &user.MongoDBRepository{Users: map[string]user.User{}}
	userService := &user.Service{Repo: userRepo}
	userHandler := &user.Handler{Svc: userService}
	r := mux.NewRouter()
	r.HandleFunc("/", userHandler.Home)
	r.HandleFunc("/create", userHandler.Create)
	r.HandleFunc("/register", userHandler.Register)
	// r.HandleFunc("/articles", ArticlesHandler)
	// http.Handle("/", r)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println(err)
	}
}
