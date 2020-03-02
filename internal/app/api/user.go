package api

import "github.com/Januadrym/seennit/internal/app/user"

func newUserService(policySvc user.PolicyService, notiSvc user.NotiService) (*user.Service, error) {
	s, err := dialMongo()
	if err != nil {
		return nil, err
	}
	repo := user.NewMongoDBRepo(s)
	// conf := user.LoadConfigFromEnv()
	return user.NewService(repo, policySvc, notiSvc), nil
}

func newUserHandler(svc *user.Service) *user.Handler {
	return user.NewHandler(svc)
}
