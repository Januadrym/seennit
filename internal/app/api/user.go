package api

import "github.com/seennit/internal/app/user"

func newUserService() (*user.Service, error) {
	s, err := dialMongo()
	if err != nil {
		return nil, err
	}
	repo := user.NewMongoDBRepo(s)
	return user.NewService(repo), nil
}

func newUserHandler(svc *user.Service) *user.Handler {
	return user.NewHandler(svc)
}
