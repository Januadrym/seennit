package comment

import (
	"context"
	"time"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/validator"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	RepoProvider interface {
		FindAll(ctx context.Context) ([]*types.Comment, error)
		FindCommentPost(ctx context.Context, idPost string) ([]*types.Comment, error)
		Create(ctx context.Context, req *types.Comment) error
		Update(ctx context.Context, id string, c string) error
		DeleteByID(ctx context.Context, id string) error
	}
	PolicyService interface {
		AddPolicy(ctx context.Context, req types.Policy) error
		Validate(ctx context.Context, obj string, act string) error
	}

	PostService interface {
		FindByID(ctx context.Context, id string) (*types.Post, error)
	}

	Service struct {
		Repo    RepoProvider
		Policy  PolicyService
		PostSvc PostService
	}
)

func NewService(repo RepoProvider, policy PolicyService, postsvc PostService) *Service {
	return &Service{
		Repo:    repo,
		Policy:  policy,
		PostSvc: postsvc,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]*types.Comment, error) {
	return s.Repo.FindAll(ctx)
}

func (s *Service) Create(ctx context.Context, req *types.Comment, idPost string) (*types.Comment, error) {
	if err := validator.Validate(req); err != nil {
		logrus.Errorf("invalid comment, err: %v", err)
		return nil, err
	}
	pt, err := s.PostSvc.FindByID(ctx, idPost)
	if err != nil {
		logrus.Errorf("failed to find post, err: %v", err)
		return nil, err
	}
	if pt.Status == types.StatusArchived {
		return nil, status.Post().Archived
	}
	thisComment := &types.Comment{
		ID:        uuid.New().String(),
		Content:   req.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		PostID:    idPost,
	}
	user := auth.FromContext(ctx)
	if user != nil {
		thisComment.CreatedByID = user.UserID
		thisComment.CreatedByName = user.GetName()
	}
	if err := validator.Validate(thisComment); err != nil {
		return nil, err
	}
	if err := s.Repo.Create(ctx, thisComment); err != nil {
		logrus.WithContext(ctx).Errorf("failed to add comment, err: %v", err)
		return nil, err
	}
	// make owner
	if err := s.Policy.AddPolicy(auth.NewAdminContext(ctx), types.Policy{
		Subject: user.UserID,
		Object:  thisComment.ID,
		Action:  types.PolicyActionAny,
		Effect:  types.PolicyEffectAllow,
	}); err != nil {
		return nil, err
	}

	return thisComment, nil
}

func (s *Service) GetAllComments(ctx context.Context, idPost string) ([]*types.Comment, error) {
	list, err := s.Repo.FindCommentPost(ctx, idPost)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find comments, err: %v", err)
		return nil, err
	}
	return list, nil
}

func (s *Service) Update(ctx context.Context, id string, c *types.Comment) error {
	if err := s.Policy.Validate(ctx, id, types.PolicyActionAny); err != nil {
		logrus.Errorf("unauthorized, not owner, err: %v", err)
		return err
	}

	return s.Repo.Update(ctx, id, c.Content)
}

func (s *Service) DeleteByID(ctx context.Context, id string) error {
	if err := s.Policy.Validate(ctx, id, types.PolicyActionAny); err != nil {
		logrus.Errorf("unauthorized, not owner, err: %v", err)
		return err
	}

	return s.Repo.DeleteByID(ctx, id)
}
