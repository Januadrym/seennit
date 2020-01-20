package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type repoProvider interface {
	Create(context.Context, User) (string, error)
	Home(context.Context, User) (string, error)
	FindUserByMail(ctx context.Context, email string) (*User, error)
	Insert(context.Context, User) error
}

type Service struct {
	Repo repoProvider
}

func (s *Service) Create(ctx context.Context, user User) (string, error) {
	id, err := s.Repo.Create(ctx, user)
	if err != nil {
		fmt.Println("handle error")
		return "", err
	}
	return id, nil
}

func (s *Service) Home(ctx context.Context, user User) (string, error) {
	name, err := s.Repo.Home(ctx, user)
	if err != nil {
		fmt.Println("handle this error")
		return "", err
	}
	return name, nil
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (string, error) {
	userDB, err := s.Repo.FindUserByMail(ctx, req.Email)
	if err != nil && err != ErrUserNotFound {
		logrus.Errorf("fail to find user: %v", err)
		return "", err
	}
	if userDB != nil {
		return "", ErrUserAlreadyExist
	}

	pword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("fail to gen password: &v", err)
		return "", fmt.Errorf("fail to register")
	}

	user := User{
		ID:        uuid.New().String(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Gender:    req.Gender,
		Password:  string(pword),
		CreatedAt: time.Now(),
	}

	if err := s.Repo.Insert(ctx, user); err != nil {
		logrus.Errorf("fail to insert: &v", err)
		return "", fmt.Errorf("fail to register: %v", err)
	}
	return user.ID, nil
}
