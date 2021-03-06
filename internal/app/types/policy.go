package types

type (
	Policy struct {
		Subject string `json:"subject" validate:"required"`
		Object  string `json:"object" validate:"required"`
		Action  string `json:"action" validate:"required"`
		Effect  string `json:"effect" validate:"oneof=* allow deny"`
	}
)

// Policy any
const (
	PolicyObjectAny  = "*"
	PolicyActionAny  = "*"
	PolicySubjectAny = "*"
)

// Policy effects
const (
	PolicyEffectAllow = "allow"
	PolicyEffectDeny  = "deny"
)

// Policy action
const (
	PolicyActionDeny = "deny"
)

// Policy of community
const (
	PolicyActionCommunity     = "community"
	PolicyActionCommunityRead = "community:read"
)
