package community

import (
	"context"

	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/jwt"

	"github.com/sirupsen/logrus"
)

type (
	repoProvider interface {
		FindCommunityByID(ctx context.Context, cID string) (*types.Community, error)
		FindAll(context.Context) ([]*types.Community, error)
	}

	Service struct {
		Jwt  jwt.SignVerifier
		Repo repoProvider
	}
)

func NewService(repo repoProvider, jwtSigner jwt.SignVerifier) *Service {
	return &Service{
		Repo: repo,
		Jwt:  jwtSigner,
	}
}

func (s *Service) SearchCommunity(ctx context.Context, req *types.Community) (*types.Community, error) {
	com, err := s.Repo.FindCommunityByID(ctx, req.CommunityID)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	return com, nil
}
