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
		Create(ctx context.Context, com *types.Community) error
		FindCommunityByID(ctx context.Context, cID string) (*types.Community, error)
		FindAllCom(context.Context) ([]*types.Community, error)
		FindCommunityByName(ctx context.Context, cName string) (*types.Community, error)
		ChangeStatus(ctx context.Context, id string, status types.Status) error
		EnrollUser(ctx context.Context, idUser string, idCom string) error
		CheckUserEnrolled(ctx context.Context, idUser string, idCom string) (string, error)
		UpdateInfo(ctx context.Context, idCom string, comm *types.Community) error
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

func (s *Service) CreateCommunity(ctx context.Context, cm *types.Community) (*types.Community, error) {
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
	comm := &types.Community{
		ID:            uuid.New().String(),
		Name:          cm.Name,
		BannerURL:     cm.BannerURL,
		Description:   cm.Description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		CreatedByID:   user.UserID,
		CreatedByName: user.GetName(),
	}
	if err := validator.Validate(comm); err != nil {
		return nil, err
	}
	if err := s.Repo.Create(ctx, comm); err != nil {
		logrus.Errorf("fail to insert: %v", err)
		return nil, fmt.Errorf("fail to create community: %v", err)
	}

	// make owner of community
	if err := s.policy.AddPolicy(auth.NewAdminContext(ctx), types.Policy{
		Subject: user.UserID,
		Object:  comm.ID,
		Action:  types.PolicyActionAny,
		Effect:  types.PolicyEffectAllow,
	}); err != nil {
		return nil, err
	}
	// enroll self
	if err := s.Repo.EnrollUser(ctx, user.UserID, comm.ID); err != nil {
		return nil, err
	}

	return comm, nil
}

// // check role
// if err := s.policy.Validate(ctx, com.ID, types.PolicyActionAny); err != nil {
// 	logrus.Info("check role info: ", err)
// 	return nil, err
// }

func (s *Service) SearchCommunity(ctx context.Context, name string) (*types.Community, error) {
	com, err := s.Repo.FindCommunityByName(ctx, name)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	if com.Status == types.StatusPrivate {
		return nil, fmt.Errorf("Community not found")
	}
	return com, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*types.Community, error) {
	com, err := s.Repo.FindAllCom(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find communities, err: %v", err)
		return nil, err
	}
	return com, nil
}

// TODO-later: status - community don't get deleted, only hidden or archive
// ATM just delete com for simple usage
func (s *Service) PrivateCommunity(ctx context.Context, Com string) error {
	com, err := s.SearchCommunity(ctx, Com)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return err
	}
	if err := s.policy.Validate(ctx, com.ID, types.PolicyActionAny); err != nil {
		logrus.Errorf("unauthorized, not owner, err: %v", err)
		return err
	}
	if err := s.Repo.ChangeStatus(ctx, com.ID, types.StatusPrivate); err != nil {
		return nil
	}
	return nil
}

func (s *Service) GetCommunity(ctx context.Context, name string) (*types.Community, error) {
	com, err := s.Repo.FindCommunityByName(ctx, name)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	if com.Status == types.StatusPrivate {
		return nil, fmt.Errorf("Community not found")
	}
	return com, nil
}

func (s *Service) EnrollUser(ctx context.Context, idCom string) error {
	user := auth.FromContext(ctx)
	id, err := s.Repo.CheckUserEnrolled(ctx, user.UserID, idCom)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Errorf("failed to check user in community, err: %v", err)
		return err
	}
	if id != "" {
		logrus.WithContext(ctx).Errorf("user already enrolled")
		return status.Community().UserEnrolled
	}
	if err := s.Repo.EnrollUser(ctx, user.UserID, idCom); err != nil {
		logrus.WithContext(ctx).Errorf("failed to enroll user, err: %v", err)
	}
	return nil

}

func (s *Service) UpdateInfo(ctx context.Context, idCom string, comm *types.Community) error {
	if err := s.policy.Validate(ctx, idCom, types.PolicyActionAny); err != nil {
		logrus.Errorf("unauthorized, not owner, err: %v", err)
		return err
	}
	return s.Repo.UpdateInfo(ctx, idCom, comm)
}

// post
// func (s *Service) AddPost(ctx context.Context, idPost string, idCom string) error {
// 	return s.Repo.AddPost(ctx, idPost, idCom)
// }

// func (s *Service) GetAllPost(ctx context.Context, idCom string) ([]string, error) {
// 	return s.Repo.GetAllPost(ctx, idCom)
// }

// func (s *Service) CheckContainPost(ctx context.Context, nameCom string, idPost string) error {
// 	in, err := s.Repo.CheckContainPost(ctx, nameCom, idPost)
// 	if err != nil && !db.IsErrNotFound(err) {
// 		logrus.WithContext(ctx).Errorf("failed to check contained post, err: %v", err)
// 		return err
// 	}
// 	if in {
// 		return nil
// 	}
// 	logrus.WithContext(ctx).Errorf("post is not contain in this community %v, err: %v", nameCom, err)
// 	return err
// }
