package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/config/env"
	"github.com/Januadrym/seennit/internal/pkg/db"
	"github.com/Januadrym/seennit/internal/pkg/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type (
	repoProvider interface {
		FindUserByMail(ctx context.Context, email string) (*types.User, error)
		Insert(ctx context.Context, user *types.User) error
		FindAll(context.Context) ([]*types.User, error)
		Delete(ctx context.Context) error
		UpdateInfo(ctx context.Context, userID string, user *types.User) error
	}

	Service struct {
		Repo repoProvider
		Conf Config
		Jwt  jwt.SignVerifier
	}

	Config struct {
		ResetPasswordTokenLifetime time.Duration `envconfig:"USER_RESET_PASSWORD_TOKEN_LIFE_TIME" default:"15m"`
	}
)

func NewService(conf Config, repo repoProvider, jwtSigner jwt.SignVerifier) *Service {
	return &Service{
		Repo: repo,
		Conf: conf,
		Jwt:  jwtSigner,
	}
}

func LoadConfigFromEnv() Config {
	var conf Config
	env.Load(&conf)
	return conf
}

func (s *Service) SearchUser(ctx context.Context, req *types.User) (*types.User, error) {
	usr, err := s.Repo.FindUserByMail(ctx, req.Email)
	if err != nil {
		fmt.Println("handle error")
		return nil, err
	}
	return usr, nil

}

func (s *Service) Register(ctx context.Context, req *types.RegisterRequest) (*types.User, error) {
	userDB, err := s.Repo.FindUserByMail(ctx, req.Email)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Errorf("failed to check user by email, err: %v", err)
		return nil, fmt.Errorf("failed to check user by email, err: %v", err)
	}

	if userDB != nil {
		logrus.WithContext(ctx).Errorf("email exist!")
		return nil, status.User().DuplicatedEmail
	}

	pword, err := s.generatePassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	user := &types.User{
		UserID:    uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  pword,
		CreatedAt: time.Now(),
	}

	if err := s.Repo.Insert(ctx, user); err != nil {
		logrus.Errorf("fail to insert: &v", err)
		return nil, fmt.Errorf("fail to register: %v", err)
	}
	return user.Strip(), nil
}

func (s *Service) FindAll(ctx context.Context) ([]*types.User, error) {
	users, err := s.Repo.FindAll(ctx)
	info := make([]*types.User, 0)
	for _, usr := range users {
		info = append(info, &types.User{
			UserID:    usr.UserID,
			Email:     usr.Email,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
		})
	}
	return info, err
}

func (s *Service) DeleteAll(ctx context.Context) error {
	if err := s.Repo.Delete(ctx); err != nil {
		return nil
	}
	return nil
}

func (s *Service) Update(ctx context.Context, userID string, user *types.User) error {
	err := s.Repo.UpdateInfo(ctx, userID, user)
	return err
}

func (s *Service) Auth(ctx context.Context, email, password string) (*types.User, error) {
	user, err := s.Repo.FindUserByMail(ctx, email)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Errorf("failed to check existing user by email, err: %v", err)
		return nil, errors.New("internal error")
	}
	if db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Debugf("user not found, email: %s", email)
		return nil, errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.WithContext(ctx).Error("invalid password")
		return nil, errors.New("invalid password")
	}
	return user.Strip(), nil
}

func (s *Service) generatePassword(pass string) (string, error) {
	rs, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password: %w", err)
	}
	return string(rs), nil
}
