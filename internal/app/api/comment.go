package api

import "github.com/Januadrym/seennit/internal/app/comment"

func newCommentService(policy comment.PolicyService, postsvc comment.PostService) (*comment.Service, error) {
	s, err := dialMongo()
	if err != nil {
		return nil, err
	}
	repo := comment.NewMongoDBRepo(s)
	return comment.NewService(repo, policy, postsvc), nil
}
func newCommentHandler(svc *comment.Service) *comment.Handler {
	return comment.NewHander(svc)
}
