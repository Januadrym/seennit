package api

import (
	"github.com/Januadrym/seennit/internal/app/policy"
	"github.com/Januadrym/seennit/internal/pkg/config/env"
)

func newPolicyService() (*policy.Service, error) {
	var conf policy.CasbinConfig
	env.LoadWithPrefix("CASBIN", &conf)
	enforcer := policy.NewMongoDBCasbinEnforcer(conf)
	return policy.New(enforcer)
}
