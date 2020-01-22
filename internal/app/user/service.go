package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"vnmquan.com/seennit/internal/app/types"
)

type (
	repoProvider interface {
		FindUserByMail(ctx context.Context, email string) (*types.User, error)
		Insert(ctx context.Context, user *types.User) error
		FindAll(context.Context) ([]*types.User, error)
		Delete(ctx context.Context) error
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

	if err != nil && err != ErrUserNotFound {
		logrus.Errorf("fail to find user: %v", err)
		return nil, err
	}

	if userDB != nil {
		logrus.Errorf("user exist: %v", err)
		return nil, ErrUserAlreadyExist
	}

	pword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("fail to gen password: &v", err)
		return nil, fmt.Errorf("fail to register")
	}

	user := &types.User{
		ID:        uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  string(pword),
		CreatedAt: time.Now(),
	}

	if err := s.Repo.Insert(ctx, user); err != nil {
		logrus.Errorf("fail to insert: &v", err)
		return nil, fmt.Errorf("fail to register: %v", err)
	}
	return user, nil
}

func (s *Service) FindAll(ctx context.Context) ([]*types.User, error) {
	users, err := s.Repo.FindAll(ctx)
	info := make([]*types.User, 0)
	for _, usr := range users {
		info = append(info, &types.User{
			ID:        usr.ID,
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
