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

// Policy of user
const (
	PolicyObjectUser         = "user"
	PolicyActionUserReadList = "solution:read-list"
)

// Policy of community
const (
	PolicyObjectCommunity       = "community"
	PolicyActionCommunityCreate = "community:create"
	PolicyActionCommunityDelete = "community:delete"
	PolicyActionCommunityUpdate = "community:update"
	PolicyActionCommunityRead   = "community:read"
)

// Policy of post
const (
	PolicyObjectPost       = "post"
	PolicyActionPostCreate = "post:create"
	PolicyActionPostDelete = "post:delete"
	PolicyActionPostUpdate = "post:update"
	PolicyActionPostRead   = "post:read"
)
