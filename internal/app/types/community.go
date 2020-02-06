package types

import "time"

type (
	Community struct {
		CommunityID   string    `json:"community_ID,omitempty" bson:"community_ID,omitempty"`
		CommunityName string    `json:"community_name,omitempty" bson:"community_name,omitempty"`
		Description   string    `json:"description,omitempty" bson:"description,omitempty"`
		BannerURL     string    `json:"banner_URL,omitempty" bson:"banner_URL,omitempty"`
		CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdatedAt     time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	}

	CommunityRequest struct {
		CommunityName string `json:"community_name,omitempty" bson:"community_name,omitempty"`
		Description   string `json:"description,omitempty" bson:"description,omitempty"`
		BannerURL     string `json:"banner_URL,omitempty" bson:"banner_URL,omitempty"`
	}
)
