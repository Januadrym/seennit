package api

import "github.com/Januadrym/seennit/internal/app/notification"

func newNotificationService() (*notification.Service, error) {
	s, err := dialMongo()
	if err != nil {
		return nil, err
	}
	repo := notification.NewMongoDBRepo(s)
	return notification.NewService(repo), nil
}
