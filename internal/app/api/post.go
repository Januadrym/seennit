package api

import (
	"github.com/Januadrym/seennit/internal/app/post"
)

func newPostService(policy post.PolicyService, cmtService post.CommentService) (*post.Service, error) {
	s, err := dialMongo()
	if err != nil {
		return nil, err
	}
	repo := post.NewMongoDBRepo(s)
	return post.NewService(repo, policy, cmtService), nil
}
func newPostHandler(svc *post.Service) *post.Handler {
	return post.NewHandler(svc)
}
