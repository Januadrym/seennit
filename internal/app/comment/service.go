package comment

import (
	"context"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/validator"
	"github.com/sirupsen/logrus"
)

type (
	RepoProvider interface {
		FindAll(ctx context.Context) ([]*types.Comment, error)
	}
	PolicyService interface {
	}
	Service struct {
		Repo   RepoProvider
		Policy PolicyService
	}
)

func NewService(repo RepoProvider, policy PolicyService) *Service {
	return &Service{
		Repo:   repo,
		Policy: policy,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]*types.Comment, error) {
	return s.Repo.FindAll(ctx)
}

func (s *Service) Create(ctx context.Context, cm *types.Comment) error {
	if err := validator.Validate(cm); err != nil {
		logrus.Errorf("invalid comment, err: %v", err)
		return err
	}
	user := auth.FromContext(ctx)
	if user != nil {
		cm.CreatedByID = user.UserID
		cm.CreatedByName = user.GetName()
	}
	// post_id for level 0 comment (not reply)
	//
	// todo
	return nil
}
