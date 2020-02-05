package types

import "time"

type (
	Community struct {
		CommunityID   string    `json:"communityID,omitempty" bson:"communityID,omitempty"`
		CommunityName string    `json:"communityName,omitempty" bson:"communityName,omitempty"`
		Description   string    `json:"description,omitempty" bson:"description,omitempty"`
		BannerURL     string    `json:"bannerURL,omitempty" bson:"bannerURL,omitempty"`
		CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdatedAt     time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	}

	CreateCommunityRequest struct {
		CommunityName string `json:"communityName,omitempty" bson:"communityName,omitempty"`
		Description   string `json:"description,omitempty" bson:"description,omitempty"`
		BannerURL     string `json:"bannerURL,omitempty" bson:"bannerURL,omitempty"`
	}
)
