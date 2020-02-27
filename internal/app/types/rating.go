package types

import "time"

type (
	RatingType       string
	RatingTargetType string
	Rating           struct {
		ID string `json:"id" bson:"id"`

		TargetID   string           `json:"target_id,omitempty" bson:"target_id" validate:"required"`
		TargetType RatingTargetType `json:"target_type,omitempty" bson:"target_type" validate:"required,oneof=post comment"`
		Type       RatingType       `json:"type" bson:"type" validate:"required,oneof=upvote downvote"`

		CreatedByID   string `json:"created_by_id" bson:"created_by_id"`
		CreatedByName string `json:"created_by_name" bson:"created_by_name"`

		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	}
)

const (
	TypeUpVote   RatingType = "upvote"
	TypeDownVote RatingType = "downvote"

	TargetTypePost    RatingTargetType = "post"
	TargetTypeComment RatingTargetType = "comment"
)
