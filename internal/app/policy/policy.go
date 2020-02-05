package policy

import (
	"context"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/db/mongodb"

	"github.com/casbin/casbin"
	"github.com/sirupsen/logrus"
)

type (
	CasbinConfig struct {
		MongoDB    mongodb.Config
		CongifPath string `envconfig:"CONFIG_PATH" default:"configs/casbin.conf"`
	}
	Service struct {
		enforcer *casbin.Enforcer
	}
)

// New return a new instance of policy service
func New(enforcer *casbin.Enforcer) (*Service, error) {
	enforcer.EnableAutoSave(true)
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}
	return &Service{
		enforcer: enforcer,
	}, nil
}

func (s *Service) addPolicy(ctx context.Context, p types.Policy) error {
	_, err := s.enforcer.AddPolicySafe(p.Subject, p.Object, p.Action, p.Effect)
	return err
}

func (s *Service) addGroupingPolicy(ctx context.Context, sub string, group string) error {
	_, err := s.enforcer.AddGroupingPolicySafe(sub, group)
	return err
}

// IsAllowed check if the sub is allowed to do the action on that object.
func (s *Service) IsAllowed(ctx context.Context, sub string, obj string, act string) bool {
	ok, err := s.enforcer.EnforceSafe(sub, obj, act)
	return err == nil && ok
}

// Validate validate if the current user is allowed to do the action on the object.
func (s *Service) Validate(ctx context.Context, obj string, act string) error {
	if auth.IsAdminContext(ctx) {
		return nil
	}
	sub := types.PolicySubjectAny
	user := auth.FromContext(ctx)
	if user != nil {
		sub = user.UserID
	}
	if !s.IsAllowed(ctx, sub, obj, act) {
		logrus.WithContext(ctx).WithFields(logrus.Fields{"sub": sub, "action": act, "obj": obj}).Errorf("the user is not authorized to do the action")
		return status.Policy().Unauthorized
	}
	return nil
}
