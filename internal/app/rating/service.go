package rating

import (
	"context"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/validator"
	"github.com/sirupsen/logrus"
)

type (
	RepoProvider interface {
		Upsert(ctx context.Context, rating *types.Rating) (bool, error)
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

func (s *Service) Create(ctx context.Context, rating *types.Rating) error {
	if err := validator.Validate(rating); err != nil {
		return status.Gen().BadRequest
	}
	user := auth.FromContext(ctx)
	if user != nil {
		rating.CreatedByID = user.UserID
		rating.CreatedByName = user.GetName()
	}
	isnew, er := s.Repo.Upsert(ctx, rating)
	if er != nil {
		logrus.WithContext(ctx).Errorf("failed to create a vote, err: %v", er)
		return er
	}
	logrus.Info("is it new? ", isnew)
	return nil
}
