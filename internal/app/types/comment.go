package types

import "time"

type (
	Comment struct {
		ID        string `json:"id,omitempty" bson:"id"`
		PostID    string `json:"post_id,omitempty" bson:"post_id" validate:"required"`
		Content   string `json:"content,omitempty" bson:"content" validate:"required"`
		ReplyToID string `json:"reply_to_id,omitempty" bson:"reply_to_id"`
		Level     int    `json:"level" bson:"level"`

		Upvote   int64 `json:"upvote" bson:"upvote"`
		Downvote int64 `json:"downvote" bson:"downvote"`

		CreatedByName string    `json:"created_by_name,omitempty" bson:"created_by_name"`
		CreatedByID   string    `json:"created_by_id,omitempty" bson:"created_by_id"`
		CreatedAt     time.Time `json:"created_at,omitempty" bson:"created_at"`
		UpdatedAt     time.Time `json:"modified_at,omitempty" bson:"updated_at"`
	}
)
