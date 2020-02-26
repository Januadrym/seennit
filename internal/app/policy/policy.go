package policy

import (
	"context"

	"github.com/Januadrym/seennit/internal/app/auth"
	"github.com/Januadrym/seennit/internal/app/status"
	"github.com/Januadrym/seennit/internal/app/types"
	"github.com/Januadrym/seennit/internal/pkg/db/mongodb"
	"github.com/Januadrym/seennit/internal/pkg/validator"

	"github.com/casbin/casbin"
	mongodbadapter "github.com/casbin/mongodb-adapter"
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

func NewMongoDBCasbinEnforcer(conf CasbinConfig) *casbin.Enforcer {
	dialInfo := conf.MongoDB.DialInfo()
	adapter := mongodbadapter.NewAdapterWithDialInfo(dialInfo)
	enforcer := casbin.NewEnforcer(conf.CongifPath, adapter)
	return enforcer
}

func (s *Service) addPolicy(ctx context.Context, p types.Policy) error {
	_, err := s.enforcer.AddPolicySafe(p.Subject, p.Object, p.Action, p.Effect)
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
		logrus.WithContext(ctx).WithFields(logrus.Fields{"SUBJECT": sub, "ACTION": act, "OBJECT": obj}).Errorf("This user is not authorized to do this action")
		return status.Policy().Unauthorized
	}
	return nil
}

func (s *Service) AddPolicy(ctx context.Context, req types.Policy) error {
	if err := validator.Validate(req); err != nil {
		return err
	}
	if err := s.Validate(ctx, req.Object, ActionPolicyUpdate); err != nil {
		return err
	}
	if err := s.addPolicy(ctx, types.Policy{
		Subject: req.Subject,
		Object:  req.Object,
		Action:  req.Action,
		Effect:  req.Effect,
	}); err != nil {
		logrus.WithContext(ctx).Errorf("fail to add policy, err: %v", err)
		return status.Gen().Internal
	}
	if req.Effect == types.PolicyEffectDeny {
		return nil
	}
	if _, err := s.enforcer.RemovePolicySafe(req.Subject, req.Object, req.Action, types.PolicyEffectDeny); err != nil {
		logrus.WithContext(ctx).Errorf("fail to cleaned up old policy, err: %v", err)
		return status.Gen().Internal
	}

	return nil
}

// For Community

// GetAllMods get all moderators of one community
func (s *Service) GetAllMods(ctx context.Context, id string) ([]string, error) {
	// filter by community id
	// "1" stand for object index
	plc := s.enforcer.GetFilteredPolicy(1, id)

	list := make([]string, 0)
	for _, p := range plc {
		if p[3] != types.PolicyEffectDeny && (p[2] == types.PolicyActionAny || p[2] == types.PolicyActionCommunity) {
			list = append(list, p[0])
		}
	}
	return list, nil
}

func (s *Service) RemovePolicy(ctx context.Context, idOwner, idCom string) error {
	_, err := s.enforcer.RemoveFilteredPolicySafe(0, idOwner, idCom)
	if err != nil {
		logrus.Errorf("failed to remove policy, err: %v", err)
		return status.Gen().Internal
	}
	return nil
}
