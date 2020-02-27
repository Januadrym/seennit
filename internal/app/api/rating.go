package api

import "github.com/Januadrym/seennit/internal/app/rating"

func newRatingService() (*rating.Service, error) {
	s, err := dialMongo()
	if err != nil {
		return nil, err
	}
	repo := rating.NewMongoDBRepo(s)
	return rating.NewService(repo), nil
}

func newRatingHandler(svc *rating.Service) *rating.Handler {
	return rating.NewHandler(svc)
}
