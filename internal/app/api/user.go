package api

import "github.com/Januadrym/seennit/internal/app/user"

func newUserService() (*user.Service, error) {
	s, err := dialMongo()
	if err != nil {
		return nil, err
	}
	repo := user.NewMongoDBRepo(s)
	conf := user.LoadConfigFromEnv()
	return user.NewService(conf, repo, newJWTSignVerifier()), nil
}

func newUserHandler(svc *user.Service) *user.Handler {
	return user.NewHandler(svc)
}
