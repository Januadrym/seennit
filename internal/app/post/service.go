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
		GetEntire(ctx context.Context) ([]*Post, error)
		Create(ctx context.Context, req *Post) error
		FindByID(ctx context.Context, id string) (*Post, error)
		GetAll(ctx context.Context, listID []string) ([]*Post, error)
		UpdatePost(ctx context.Context, id string, p *PostUpdateRequest) error
		ChangeStatus(ctx context.Context, id string, status Status) error
	}

	PolicyService interface {
		AddPolicy(ctx context.Context, req types.Policy) error
		Validate(ctx context.Context, obj string, act string) error
	}

	CommunityService interface {
		AddPost(ctx context.Context, idPost string, idCom string) error
		SearchCommunity(ctx context.Context, name string) (*types.Community, error)
		GetAllPost(ctx context.Context, idCom string) ([]string, error)
		CheckContainPost(ctx context.Context, nameCom string, idPost string) error
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
func (s *Service) Create(ctx context.Context, req *Post, nameComm string) (*Post, error) {
	if err := validator.Validate(req); err != nil {
		return nil, fmt.Errorf("invalid post: %v", err)
	}
	thispost := &Post{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Content:     req.Content,
		CreatedAt:   time.Now(),
		PublishDate: time.Now(),
	}

	//track who create this post
	user := auth.FromContext(ctx)
	if user != nil {
		thispost.CreatedByID = user.UserID
		thispost.CreatedByName = user.GetName()
	}
	thispost.Status = StatusPublic
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

	// add post to community
	com, err := s.community.SearchCommunity(ctx, nameComm)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	if err := s.community.AddPost(ctx, thispost.ID, com.ID); err != nil {
		logrus.WithContext(ctx).Errorf("failed to add post to community, err: %v", err)
		return nil, err
	}

	return thispost, nil
}

func (s *Service) GetAll(ctx context.Context, nameComm string) ([]*Post, error) {
	com, err := s.community.SearchCommunity(ctx, nameComm)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	list, err := s.community.GetAllPost(ctx, com.ID)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to get posts, err: %v", err)
		return nil, err
	}
	listPost, err := s.Repo.GetAll(ctx, list)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find post, err: %v", err)
		return nil, err
	}
	return listPost, nil
}

func (s *Service) UpdatePost(ctx context.Context, id string, p *PostUpdateRequest) error {
	if err := validator.Validate(p); err != nil {
		return fmt.Errorf("invalid post, err: %v", err)
	}
	if err := s.policy.Validate(ctx, id, types.PolicyActionPostUpdate); err != nil {
		return err
	}
	logrus.Info("post service ne:", p)
	return s.Repo.UpdatePost(ctx, id, p)
}

func (s *Service) FindByID(ctx context.Context, id string, nameCom string) (*Post, error) {
	p, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.community.CheckContainPost(ctx, nameCom, id); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ChangeStatus(ctx context.Context, id string, status Status) error {
	if err := s.policy.Validate(ctx, id, types.PolicyActionPostUpdate); err != nil {
		return err
	}
	return s.Repo.ChangeStatus(ctx, id, status)
}

func (s *Service) GetEntire(ctx context.Context) ([]*Post, error) {
	list, err := s.Repo.GetEntire(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("cannot get anything, err: %v", err)
		return nil, err
	}
	return list, nil
}
