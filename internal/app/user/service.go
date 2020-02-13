package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/db"
	"github.com/Januadrym/seennit/internal/pkg/validator"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type (
	repoProvider interface {
		FindUserByMail(ctx context.Context, email string) (*types.User, error)
		Create(ctx context.Context, user *types.User) error
		FindAll(context.Context) ([]*types.User, error)
		DeleteAll(ctx context.Context) error
		Delete(ctx context.Context, id string) error
		UpdateInfo(ctx context.Context, userID string, user *types.User) error
	}

	Service struct {
		Repo repoProvider
	}
)

func NewService(repo repoProvider) *Service {
	return &Service{
		Repo: repo,
	}
}

// func LoadConfigFromEnv() Config {
// 	var conf Config
// 	env.Load(&conf)
// 	return conf
// }

func (s *Service) SearchUser(ctx context.Context, req *types.User) (*types.User, error) {
	usr, err := s.Repo.FindUserByMail(ctx, req.Email)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find user by email, err: %v", err)
		return nil, err
	}
	return usr, nil
}

func (s *Service) Register(ctx context.Context, req *types.RegisterRequest) (*types.User, error) {
	if err := validator.Validate(req); err != nil {
		return nil, err
	}
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
	if err := validator.Validate(user); err != nil {
		return nil, err
	}

	if err := s.Repo.Create(ctx, user); err != nil {
		logrus.Errorf("fail to insert: %v", err)
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

func (s *Service) RemoveAll(ctx context.Context) error {
	if err := s.Repo.DeleteAll(ctx); err != nil {
		return nil
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, userID string) error {
	if err := s.Repo.Delete(ctx, userID); err != nil {
		return nil
	}
	return nil
}

func (s *Service) Update(ctx context.Context, userID string, user *types.User) error {
	if err := validator.Validate(user); err != nil {
		return err
	}

	return s.Repo.UpdateInfo(ctx, userID, user)
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
