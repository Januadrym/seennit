package api

import "github.com/Januadrym/seennit/internal/app/comment"

func newCommentService(policy comment.PolicyService) (*comment.Service, error) {
	s, err := dialMongo()
	if err != nil {
		return nil, err
	}
	repo := comment.NewMongoDBRepo(s)
	return comment.NewService(repo, policy), nil
}
func newCommentHandler(svc *comment.Service) *comment.Handler {
	return comment.NewHander(svc)
}
