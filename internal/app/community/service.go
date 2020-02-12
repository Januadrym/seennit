package community

import (
	"context"
	"fmt"
	"time"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/db"
	"github.com/Januadrym/seennit/internal/pkg/jwt"
	"github.com/Januadrym/seennit/internal/pkg/validator"
	"github.com/google/uuid"

	"github.com/sirupsen/logrus"
)

type (
	RepoProvider interface {
		Create(ctx context.Context, com *Community) error
		FindCommunityByID(ctx context.Context, cID string) (*Community, error)
		FindAllCom(context.Context) ([]*Community, error)
		FindCommunityByName(ctx context.Context, cName string) (*Community, error)
		DeleteByID(ctx context.Context, id string) error
		EnrollUser(ctx context.Context, idUser string, idCom string) error
	}

	PolicyService interface {
		AddPolicy(ctx context.Context, req types.Policy) error
		Validate(ctx context.Context, obj string, act string) error
	}

	Service struct {
		Jwt    jwt.SignVerifier
		Repo   RepoProvider
		policy PolicyService
	}
)

func NewService(repo RepoProvider, jwtSigner jwt.SignVerifier, policySvc PolicyService) *Service {
	return &Service{
		Repo:   repo,
		Jwt:    jwtSigner,
		policy: policySvc,
	}
}

func (s *Service) CreateCommunity(ctx context.Context, cm *Community) (*Community, error) {
	if err := validator.Validate(cm); err != nil {
		return nil, err
	}
	user := auth.FromContext(ctx)
	logrus.Printf("user create this: %v", user.FirstName)

	comDB, err := s.Repo.FindCommunityByName(ctx, cm.Name)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Errorf("failed to check community's name, err: %v", err)
		return nil, fmt.Errorf("failed to check community's name, err: %v", err)
	}
	if comDB != nil {
		logrus.WithContext(ctx).Errorf("name taken!")
		return nil, status.Community().NameTaken
	}
	comm := &Community{
		ID:          uuid.New().String(),
		Name:        cm.Name,
		BannerURL:   cm.BannerURL,
		Description: cm.Description,
		Users:       cm.Users,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := validator.Validate(comm); err != nil {
		return nil, err
	}
	if err := s.Repo.Create(ctx, comm); err != nil {
		logrus.Errorf("fail to insert: %v", err)
		return nil, fmt.Errorf("fail to create community: %v", err)
	}

	// make owner
	if err := s.policy.AddPolicy(auth.NewAdminContext(ctx), types.Policy{
		Subject: user.UserID,
		Object:  comm.ID,
		Action:  types.PolicyActionAny,
		Effect:  types.PolicyEffectAllow,
	}); err != nil {
		return nil, err
	}

	return comm, nil
}

// // check role
// if err := s.policy.Validate(ctx, com.ID, types.PolicyActionAny); err != nil {
// 	logrus.Info("check role info: ", err)
// 	return nil, err
// }

func (s *Service) SearchCommunity(ctx context.Context, name string) (*Community, error) {
	com, err := s.Repo.FindCommunityByName(ctx, name)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	return com, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*Community, error) {
	com, err := s.Repo.FindAllCom(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find communities, err: %v", err)
		return nil, err
	}
	return com, nil
}

// TODO-later: status - community don't get deleted, only hidden or archive
// ATM just delete com for simple usage
func (s *Service) DeleteCommunity(ctx context.Context, comID string) error {
	if err := s.policy.Validate(ctx, comID, types.PolicyActionAny); err != nil {
		logrus.Info("check role info:", err)
		return err
	}
	if err := s.Repo.DeleteByID(ctx, comID); err != nil {
		return nil
	}
	return nil
}

func (s *Service) GetCommunity(ctx context.Context, name string) (*Community, error) {
	com, err := s.Repo.FindCommunityByName(ctx, name)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
	}
	return com, nil
}

func (s *Service) EnrollUser(ctx context.Context, idCom string) error {
	user := auth.FromContext(ctx)
	err := s.Repo.EnrollUser(ctx, user.UserID, idCom)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to enroll user, err: %v", err)
	}
	return nil

}

// TODO update
