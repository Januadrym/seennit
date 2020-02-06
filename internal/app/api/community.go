package api

import "github.com/Januadrym/seennit/internal/app/community"

func newCommunityService() (*community.Service, error) {
	s, err := dialMongo()
	if err != nil {
		return nil, err
	}
	repo := community.NewMongoDBRepo(s)
	return community.NewService(repo, newJWTSignVerifier()), nil
}

func newCommunityHandler(svc *community.Service) *community.Handler {
	return community.NewHandler(svc)
}
