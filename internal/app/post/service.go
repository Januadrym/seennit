package post

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
		Create(ctx context.Context, req *types.Post) error
		GetAllPost(ctx context.Context, idCom string) ([]*types.Post, error)

		CheckPostBelongTo(ctx context.Context, idCom string, idPost string) (*types.Post, error)
		FindByID(ctx context.Context, idPost string) (*types.Post, error)
		UpdatePost(ctx context.Context, id string, p *types.PostUpdateRequest) error
		ChangeStatus(ctx context.Context, id string, status types.Status) error

		GetEntire(ctx context.Context) ([]*types.Post, error)
	}

	PolicyService interface {
		AddPolicy(ctx context.Context, req types.Policy) error
		Validate(ctx context.Context, obj string, act string) error
	}

	CommentService interface {
		Create(ctx context.Context, req *types.Comment, idPost string) (*types.Comment, error)
		GetCommentsPost(ctx context.Context, idPost string) ([]*types.Comment, error)
	}

	Service struct {
		Repo       RepoProvider
		policy     PolicyService
		cmtService CommentService
	}
)

func NewService(repo RepoProvider, policy PolicyService, cmtService CommentService) *Service {
	return &Service{
		Repo:       repo,
		policy:     policy,
		cmtService: cmtService,
	}
}
func (s *Service) Create(ctx context.Context, req *types.Post, idCom string) (*types.Post, error) {
	thispost := &types.Post{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Content:     req.Content,
		CreatedAt:   time.Now(),
		PublishDate: time.Now(),
	}
	thispost.CommunityID = idCom

	// track who create this post
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
		return nil, status.Gen().Internal
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

func (s *Service) GetAll(ctx context.Context, idCom string) ([]*types.Post, error) {
	list, err := s.Repo.GetAllPost(ctx, idCom)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to get posts, err: %v", err)
		return nil, err
	}
	return list, nil
}

func (s *Service) CheckPostBelongTo(ctx context.Context, idCom string, idPost string) (*types.Post, error) {
	p, err := s.Repo.CheckPostBelongTo(ctx, idCom, idPost)
	if err != nil {
		return nil, err
	}
	if p.Status == types.StatusDelete {
		return nil, status.Post().NotFound
	}
	return p, nil
}

func (s *Service) FindByID(ctx context.Context, idPost string) (*types.Post, error) {
	p, err := s.Repo.FindByID(ctx, idPost)
	if err != nil {
		return nil, err
	}
	if p.Status == types.StatusDelete {
		return nil, status.Post().NotFound
	}
	return p, nil
}

func (s *Service) UpdatePost(ctx context.Context, idPost string, p *types.PostUpdateRequest) error {
	if err := validator.Validate(p); err != nil {
		return status.Gen().BadRequest
	}
	if err := s.policy.Validate(ctx, idPost, types.PolicyActionAny); err != nil {
		return err
	}
	post, err := s.Repo.FindByID(ctx, idPost)
	if err != nil {
		return err
	}
	if post.Status == types.StatusArchived {
		return status.Post().Archived
	}
	return s.Repo.UpdatePost(ctx, idPost, p)
}

func (s *Service) ChangeStatus(ctx context.Context, idPost string, stat types.Status) error {
	if err := s.policy.Validate(ctx, idPost, types.PolicyActionAny); err != nil {
		return err
	}
	_, err := s.FindByID(ctx, idPost)
	if err != nil {
		return status.Gen().NotFound
	}
	return s.Repo.ChangeStatus(ctx, idPost, stat)
}

// Get entire
// get all post in all community
// homepage
// TODO: ranking by vote

func (s *Service) GetEntire(ctx context.Context) ([]*types.Post, error) {
	list, err := s.Repo.GetEntire(ctx)
	if err != nil {
		logrus.WithContext(ctx).Errorf("cannot get anything, err: %v", err)
		return nil, err
	}
	return list, nil
}

// Comments
func (s *Service) CreateComment(ctx context.Context, req *types.Comment, idPost string) (*types.Comment, error) {
	pt, err := s.FindByID(ctx, idPost)
	if err != nil {
		logrus.Errorf("failed to find post, err: %v", err)
		return nil, err
	}
	if pt.Status == types.StatusArchived {
		return nil, status.Post().Archived
	}
	p, err := s.cmtService.Create(ctx, req, idPost)
	if err != nil {
		logrus.Errorf("failed to create comment, err: %v", err)
		return nil, status.Gen().Internal
	}
	return p, nil
}

func (s *Service) GetCommentsPost(ctx context.Context, idPost string) ([]*types.Comment, error) {
	return s.cmtService.GetCommentsPost(ctx, idPost)
}
