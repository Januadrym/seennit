package post

import "time"

type (
	Status string
	Post   struct {
		ID          string    `json:"id" bson:"id"`
		Title       string    `json:"title" bson:"title" validate:"required"`
		Content     string    `json:"content" bson:"content" validate:"required"`
		Status      Status    `json:"status" bson:"status"`
		PublishDate time.Time `json:"publish_date" bson:"publish_date"`

		Tags []string `json:"tags" bson:"tags"` // TODO

		Views    int64 `json:"views" bson:"views"`
		Comments int64 `json:"comments" bson:"comments"`

		Upvote   int64 `json:"upvote" bson:"upvote"`
		Downvote int64 `json:"downvote" bson:"downvote"`

		CreatedByID   string `json:"created_by_id" bson:"created_by_id"`
		CreatedByName string `json:"created_by_name" bson:"created_by_name"`

		CreatedAt time.Time `json:"created_at" bson:"created_at"`
		UpdatedAt time.Time `json:"update_at" bson:"updated_at"`
	}
)

const (
	StatusPublic   Status = "public"
	StatusArchived Status = "archived"
	StatusDraft    Status = "draft"
	StatusDelete   Status = "deleted"
)
