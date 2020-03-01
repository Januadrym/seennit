package api

import "github.com/Januadrym/seennit/internal/app/community"

func newCommunityService(policy community.PolicyService, postSvc community.PostService, userSvc community.UserService, notiSvc community.NotifiService) (*community.Service, error) {
	s, err := dialMongo()
	if err != nil {
		return nil, err
	}
	repo := community.NewMongoDBRepo(s)
	return community.NewService(repo, policy, postSvc, userSvc, notiSvc), nil
}

func newCommunityHandler(svc *community.Service) *community.Handler {
	return community.NewHandler(svc)
}
