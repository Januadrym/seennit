package auth

import (
	"context"
	"errors"

	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/jwt"
	"github.com/sirupsen/logrus"
)

type (
	UserAuthen interface {
		Auth(ctx context.Context, email, password string) (*types.User, error)
	}
	Service struct {
		jwtSigner     jwt.Signer
		authenticator UserAuthen
	}
)

func NewService(signer jwt.Signer, authenticator UserAuthen) *Service {
	return &Service{
		jwtSigner:     signer,
		authenticator: authenticator,
	}
}

func (s *Service) Auth(ctx context.Context, email, password string) (string, *types.User, error) {
	user, err := s.authenticator.Auth(ctx, email, password)
	if err != nil {
		logrus.WithContext(ctx).Errorf("fail to login with %s, err: %#v", email, err)
		return "", nil, errors.New("unauthorized")
	}
	token, err := s.jwtSigner.Sign(userToClaims(user, jwt.DefaultLifeTime))
	if err != nil {
		logrus.WithContext(ctx).Errorf("fail to gerenate JWT token, err: %#v", err)
		return "", nil, err
	}
	return token, user, nil
}
