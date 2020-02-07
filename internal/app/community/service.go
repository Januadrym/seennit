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
		Delete(ctx context.Context) error
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

func (s *Service) SearchCommunity(ctx context.Context, req *types.Community) (*types.Community, error) {
	com, err := s.Repo.FindCommunityByID(ctx, req.CommunityID)
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to find community, err: %v", err)
		return nil, err
	}
	return com, nil
}

func (s *Service) CreateCommunity(ctx context.Context, req *types.CommunityRequest) (*types.Community, error) {
	if err := validator.Validate(req); err != nil {
		return nil, err
	}
	user := auth.FromContext(ctx)
	logrus.Printf("user create this %v", user)
	user.Roles = append(user.Roles, "admin")

	comDB, err := s.Repo.FindCommunityByName(ctx, req.CommunityName)
	if err != nil && !db.IsErrNotFound(err) {
		logrus.WithContext(ctx).Errorf("failed to check community's name, err: %v", err)
		return nil, fmt.Errorf("failed to check community's name, err: %v", err)
	}
	if comDB != nil {
		logrus.WithContext(ctx).Errorf("name taken!")
		return nil, status.Community().NameTaken
	}
	comm := &types.Community{
		CommunityID:   uuid.New().String(),
		BannerURL:     req.BannerURL,
		CommunityName: req.CommunityName,
		CreatedAt:     time.Now(),
		Description:   req.Description,
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
		Object:  comm.CommunityID,
		Action:  types.PolicyActionAny,
		Effect:  types.PolicyEffectAllow,
	}); err != nil {
		return nil, err
	}

	return comm, nil
}
