package user

import (
	"context"
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
		Delete(ctx context.Context, id string) error
		UpdateInfo(ctx context.Context, userID string, user *types.User) error
		EnrollUser(ctx context.Context, idUser string, idCom string) error
		CheckUserEnrolled(ctx context.Context, idUser string, idCom string) (string, error)
		GetUsersCommunity(ctx context.Context, idCom string) ([]*types.User, error)
		GetMods(ctx context.Context, listID []string) ([]*types.User, error)
	}

	PolicyService interface {
		Validate(ctx context.Context, obj string, act string) error
	}

	NotiService interface {
		LoadNotiUser(ctx context.Context, userID string) ([]*types.PushNotification, error)
	}

	Service struct {
		Repo    repoProvider
		policy  PolicyService
		notiSvc NotiService
	}
)

func NewService(repo repoProvider, policySvc PolicyService, notiSvc NotiService) *Service {
	return &Service{
		Repo:    repo,
		policy:  policySvc,
		notiSvc: notiSvc,
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
		return nil, status.Gen().Internal
	}

	if userDB != nil {
		logrus.WithContext(ctx).Errorf("email exist!")
		return nil, status.User().DuplicatedEmail
	}

	pword, err := s.generatePassword(req.Password)
	if err != nil {
		return nil, status.Gen().Internal
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
		return nil, status.User().RegisterFail
	}
	return user.Strip(), nil
}

func (s *Service) FindAll(ctx context.Context) ([]*types.User, error) {
	if err := s.policy.Validate(ctx, types.PolicyObjectAny, types.PolicyActionAny); err != nil {
		return nil, err
	}
	users, err := s.Repo.FindAll(ctx)
	if err != nil {
		logrus.Errorf("failed to find users, err: %v", err)
		return nil, err
	}
	info := make([]*types.User, 0)
	for _, usr := range users {
		info = append(info, &types.User{
			UserID:    usr.UserID,
			Email:     usr.Email,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
			AvatarURL: usr.AvatarURL,
		})
	}
	return info, nil
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
		return nil, status.Gen().Internal
	}
	if db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Debugf("user not found, email: %s", email)
		return nil, status.Gen().Internal
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logrus.WithContext(ctx).Error("invalid password")
		return nil, status.Gen().Internal
	}
	return user.Strip(), nil
}

func (s *Service) generatePassword(pass string) (string, error) {
	rs, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("fail to gen password, err: %v", err)
		return "", status.Gen().Internal
	}
	return string(rs), nil
}

func (s *Service) EnrollUser(ctx context.Context, idCom, idUser string) error {
	return s.Repo.EnrollUser(ctx, idUser, idCom)
}

func (s *Service) CheckUserEnrolled(ctx context.Context, idUser string, idCom string) (string, error) {
	return s.Repo.CheckUserEnrolled(ctx, idUser, idCom)
}

func (s *Service) GetUsersCommunity(ctx context.Context, idCom string) ([]*types.User, error) {
	users, err := s.Repo.GetUsersCommunity(ctx, idCom)
	if err != nil {
		logrus.Errorf("failed to get users, err: %v", err)
		return nil, err
	}
	info := make([]*types.User, 0)
	for _, usr := range users {
		info = append(info, &types.User{
			UserID:    usr.UserID,
			Email:     usr.Email,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
			AvatarURL: usr.AvatarURL,
		})
	}
	return info, nil
}

func (s *Service) GetMods(ctx context.Context, listID []string) ([]*types.User, error) {
	listMods, err := s.Repo.GetMods(ctx, listID)
	if err != nil {
		logrus.Errorf("failed to get mods, err: %v", err)
		return nil, err
	}
	info := make([]*types.User, 0)
	for _, usr := range listMods {
		info = append(info, &types.User{
			UserID:    usr.UserID,
			Email:     usr.Email,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
			AvatarURL: usr.AvatarURL,
		})
	}
	return info, nil
}

func (s *Service) LoadNotification(ctx context.Context, ID string) ([]*types.PushNotification, error) {
	return s.notiSvc.LoadNotiUser(ctx, ID)
}
