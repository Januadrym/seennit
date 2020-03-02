package notification

import (
	"context"

	"github.com/Januadrym/seennit/internal/app/types"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	RepoProvider interface {
		Create(ctx context.Context, noti *types.PushNotification) error
		LoadNoti(ctx context.Context, userID string) ([]*types.PushNotification, error)
	}

	Service struct {
		Repo RepoProvider
	}
)

func NewService(repo RepoProvider) *Service {
	return &Service{
		Repo: repo,
	}
}

//

func (s *Service) CreateNotificaion(ctx context.Context, postID string, userIDs []string, mess string) error {
	newmsg := "New post has just created. Title: " + mess
	noti := &types.PushNotification{
		ID:      uuid.New().String(),
		Message: newmsg,
		PostID:  postID,
		UserID:  userIDs,
	}

	if err := s.Repo.Create(ctx, noti); err != nil {
		logrus.WithContext(ctx).Errorf("failed to create notification, err: %v", err)
		return err
	}
	return nil
}

func (s *Service) LoadNotiUser(ctx context.Context, userID string) ([]*types.PushNotification, error) {
	notis, err := s.Repo.LoadNoti(ctx, userID)
	if err != nil {
		logrus.Errorf("failed to load notification, err: %v", err)
		return nil, err
	}

	simpleNoti := make([]*types.PushNotification, 0)
	for _, noti := range notis {
		simpleNoti = append(simpleNoti, &types.PushNotification{
			Message: noti.Message,
			PostID:  noti.PostID,
		})
	}
	return simpleNoti, nil
}
