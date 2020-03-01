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
	noti := &types.PushNotification{
		ID:      uuid.New().String(),
		Message: mess,
		PostID:  postID,
		UserID:  userIDs,
	}

	if err := s.Repo.Create(ctx, noti); err != nil {
		logrus.WithContext(ctx).Errorf("failed to create notification, err: %v", err)
		return err
	}
	return nil
}

// func (s *Service) GetNotiUser(ctx context.Context, userID string) (*types.PushNotification, error) {
// 	return nil, nil
// }
