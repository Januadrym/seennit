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
	repoProvider interface {
		Create(ctx context.Context, post *Post) error
	}

	PolicyService interface {
		AddPolicy(ctx context.Context, req types.Policy) error
		Validate(ctx context.Context, obj string, act string) error
	}

	Service struct {
		Repo   repoProvider
		policy PolicyService
	}
)

func NewService(repo repoProvider, policy PolicyService) *Service {
	return &Service{
		Repo:   repo,
		policy: policy,
	}
}
func (s *Service) Create(ctx context.Context, post *Post) (*Post, error) {
	if err := validator.Validate(post); err != nil {
		return nil, fmt.Errorf("invalid post: %v", err)
	}
	thispost := &Post{
		ID:          uuid.New().String(),
		Title:       post.Title,
		Content:     post.Content,
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
	return thispost, nil
}

// TODO: community has posts
