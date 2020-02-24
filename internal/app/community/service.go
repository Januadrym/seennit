package community

import (
	"context"
	"time"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/db"
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
		ChangeStatus(ctx context.Context, id string, status types.CommunityStatus) error
		UpdateInfo(ctx context.Context, idCom string, comm *types.Community) error
	}

	PolicyService interface {
		AddPolicy(ctx context.Context, req types.Policy) error
		Validate(ctx context.Context, obj string, act string) error
	}

	PostService interface {
		Create(ctx context.Context, req *types.Post, idCom string) (*types.Post, error)
		GetAll(ctx context.Context, idCom string) ([]*types.Post, error)
	}

	UserService interface {
		EnrollUser(ctx context.Context, idCom, idUser string) error
		CheckUserEnrolled(ctx context.Context, idUser string, idCom string) (string, error)
	}

	Service struct {
		Repo        RepoProvider
		policy      PolicyService
		postService PostService
		userService UserService
	}
)

func NewService(repo RepoProvider, policySvc PolicyService, postSvc PostService, userSvc UserService) *Service {
	return &Service{
		Repo:        repo,
		policy:      policySvc,
		postService: postSvc,
		userService: userSvc,
	}
}

func (s *Service) CreateCommunity(ctx context.Context, cm *types.Community) (*types.Community, error) {
	if err := validator.Validate(cm); err != nil {
		return nil, err
	}
	user := auth.FromContext(ctx)
	logrus.Info("user create this: ", user.GetName())

	comDB, err := s.Repo.FindCommunityByName(ctx, cm.Name)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Errorf("failed to check community's name, err: %v", err)
		return nil, err
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

	// enroll self
	if err := s.userService.EnrollUser(ctx, comm.ID, user.UserID); err != nil {
		return nil, err
	}

	if err := s.Repo.Create(ctx, comm); err != nil {
		logrus.Errorf("fail to insert: %v", err)
		return nil, status.Community().CreateFail
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
	return comm, nil
}

func (s *Service) SearchCommunity(ctx context.Context, name string) (*types.Community, error) {
	com, err := s.Repo.FindCommunityByName(ctx, name)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	if com.Status == types.CommunityStatusPrivate {
		return nil, status.Community().NotFound
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
	if err := s.Repo.ChangeStatus(ctx, com.ID, types.CommunityStatusPrivate); err != nil {
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
	if com.Status == types.CommunityStatusPrivate {
		return nil, status.Community().NotFound
	}
	return com, nil
}

func (s *Service) UpdateInfo(ctx context.Context, idCom string, comm *types.Community) error {
	if err := s.policy.Validate(ctx, idCom, types.PolicyActionCommunityUpdate); err != nil {
		logrus.Errorf("unauthorized, not owner, err: %v", err)
		return err
	}
	return s.Repo.UpdateInfo(ctx, idCom, comm)
}

// Users

func (s *Service) EnrollUser(ctx context.Context, idCom string) error {
	user := auth.FromContext(ctx)
	id, err := s.userService.CheckUserEnrolled(ctx, user.UserID, idCom)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Errorf("failed to check user in community, err: %v", err)
		return err
	}
	if id != "" {
		logrus.WithContext(ctx).Errorf("user already enrolled")
		return status.Community().UserEnrolled
	}
	if err := s.userService.EnrollUser(ctx, idCom, user.UserID); err != nil {
		logrus.WithContext(ctx).Errorf("failed to enroll user, err: %v", err)
	}
	return nil

}

func (s *Service) PromoteUser(ctx context.Context, idUser string, idCom string) error {
	if err := s.policy.Validate(ctx, idCom, types.PolicyActionAny); err != nil {
		logrus.Errorf("unauthorized, not owner, err: %v", err)
		return err
	}
	if err := s.policy.AddPolicy(auth.NewAdminContext(ctx), types.Policy{
		Subject: idUser,
		Object:  idCom,
		Action:  types.PolicyActionCommunityUpdate,
		Effect:  types.PolicyEffectAllow,
	}); err != nil {
		return err
	}
	return nil
}

// Posts

func (s *Service) SubmitPost(ctx context.Context, nameComm string, req *types.Post) (*types.Post, error) {
	if err := validator.Validate(req); err != nil {
		return nil, status.Gen().BadRequest
	}
	// add post to community
	com, err := s.SearchCommunity(ctx, nameComm)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	thispost, err := s.postService.Create(ctx, req, com.ID)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to create post, err: %v", err)
		return nil, status.Gen().Internal
	}
	return thispost, nil
}

func (s *Service) GetAllPosts(ctx context.Context, nameComm string) ([]*types.Post, error) {
	com, err := s.SearchCommunity(ctx, nameComm)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	return s.postService.GetAll(ctx, com.ID)
}
