package types

import "time"

type (
	Status string

	Post struct {
		ID          string    `json:"id" bson:"id"`
		CommunityID string    `json:"community_id" bson:"community_id"`
		Title       string    `json:"title" bson:"title" validate:"required"`
		Content     string    `json:"content" bson:"content" validate:"required"`
		Status      Status    `json:"status" bson:"status"`
		PublishDate time.Time `json:"publish_date" bson:"publish_date"`

		Views int64 `json:"views" bson:"views"`

		Upvote   int64 `json:"upvote" bson:"upvote"`
		Downvote int64 `json:"downvote" bson:"downvote"`

		CreatedByID   string `json:"created_by_id" bson:"created_by_id"`
		CreatedByName string `json:"created_by_name" bson:"created_by_name"`

		CreatedAt time.Time `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time `json:"update_at" bson:"updated_at"`
	}

	PostUpdateRequest struct {
		Title     string    `json:"title" bson:"title" validate:"required"`
		Content   string    `json:"content" bson:"content" validate:"required"`
		UpdatedAt time.Time `json:"update_at" bson:"updated_at"`
	}
)

const (
	StatusPublic   Status = "public"
	StatusDraft    Status = "draft"
	StatusDelete   Status = "deleted"
	StatusPrivate  Status = "private"
	StatusArchived Status = "archived" // for reading purposes, cannot edit, comment on a archived post
)
