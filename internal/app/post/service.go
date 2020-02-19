package post

import (
	"context"
	"fmt"
	"time"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/validator"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	RepoProvider interface {
		GetEntire(ctx context.Context) ([]*types.Post, error)
		Create(ctx context.Context, req *types.Post) error
		FindByID(ctx context.Context, id string) (*types.Post, error)
		GetAllPost(ctx context.Context, idCom string) ([]*types.Post, error)
		UpdatePost(ctx context.Context, id string, p *types.PostUpdateRequest) error
		ChangeStatus(ctx context.Context, id string, status types.Status) error
	}

	PolicyService interface {
		AddPolicy(ctx context.Context, req types.Policy) error
		Validate(ctx context.Context, obj string, act string) error
	}

	CommunityService interface {
		SearchCommunity(ctx context.Context, name string) (*types.Community, error)
	}

	Service struct {
		Repo      RepoProvider
		policy    PolicyService
		community CommunityService
	}
)

func NewService(repo RepoProvider, policy PolicyService, community CommunityService) *Service {
	return &Service{
		Repo:      repo,
		policy:    policy,
		community: community,
	}
}
func (s *Service) Create(ctx context.Context, req *types.Post, nameComm string) (*types.Post, error) {
	if err := validator.Validate(req); err != nil {
		return nil, fmt.Errorf("invalid post: %v", err)
	}
	thispost := &types.Post{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Content:     req.Content,
		CreatedAt:   time.Now(),
		PublishDate: time.Now(),
	}
	// add post to community
	com, err := s.community.SearchCommunity(ctx, nameComm)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	thispost.CommunityID = com.ID

	//track who create this post
	user := auth.FromContext(ctx)
	if user != nil {
		thispost.CreatedByID = user.UserID
		thispost.CreatedByName = user.GetName()
	}
	thispost.Status = types.StatusPublic
	if err := validator.Validate(thispost); err != nil {
		return nil, err
	}
	if err := s.Repo.Create(ctx, thispost); err != nil {
		logrus.WithContext(ctx).Errorf("failed to create post, err: %v", err)
		return nil, fmt.Errorf("failed to insert post: %v", err)
	}

	// make owner of post
	if err := s.policy.AddPolicy(auth.NewAdminContext(ctx), types.Policy{
		Subject: user.UserID,
		Object:  thispost.ID,
		Action:  types.PolicyActionAny,
		Effect:  types.PolicyEffectAllow,
	}); err != nil {
		return nil, err
	}

	return thispost, nil
}

func (s *Service) GetAll(ctx context.Context, nameComm string) ([]*types.Post, error) {
	com, err := s.community.SearchCommunity(ctx, nameComm)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	list, err := s.Repo.GetAllPost(ctx, com.ID)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to get posts, err: %v", err)
		return nil, err
	}
	return list, nil
}

func (s *Service) UpdatePost(ctx context.Context, id string, p *types.PostUpdateRequest) error {
	if err := validator.Validate(p); err != nil {
		return fmt.Errorf("invalid post, err: %v", err)
	}
	if err := s.policy.Validate(ctx, id, types.PolicyActionPostUpdate); err != nil {
		return err
	}
	logrus.Info("post service ne:", p)
	return s.Repo.UpdatePost(ctx, id, p)
}

func (s *Service) FindByID(ctx context.Context, id string) (*types.Post, error) {
	p, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ChangeStatus(ctx context.Context, id string, status types.Status) error {
	if err := s.policy.Validate(ctx, id, types.PolicyActionPostUpdate); err != nil {
		return err
	}
	return s.Repo.ChangeStatus(ctx, id, status)
}

func (s *Service) GetEntire(ctx context.Context) ([]*types.Post, error) {
	list, err := s.Repo.GetEntire(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("cannot get anything, err: %v", err)
		return nil, err
	}
	return list, nil
}
