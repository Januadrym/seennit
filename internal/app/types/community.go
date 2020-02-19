package types

import "time"

type (
	Community struct {
		ID          string `json:"id,omitempty" bson:"id,omitempty"`
		Name        string `json:"name,omitempty" bson:"name,omitempty"`
		Description string `json:"description,omitempty" bson:"description,omitempty"`
		BannerURL   string `json:"banner_url,omitempty" bson:"banner_url,omitempty"`

		Users []string `json:"users,omitempty" bson:"users,omitempty"` // consider

		CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdatedAt     time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
		CreatedByID   string    `json:"created_by_id,omitempty" bson:"created_by_id,omitempty"`
		CreatedByName string    `json:"created_by_name,omitempty" bson:"created_by_name,omitempty"`
	}
)
