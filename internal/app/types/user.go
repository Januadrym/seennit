package types

import (
	"time"
)

type (
	RegisterRequest struct {
		FirstName string `validate:"required" json:"first_name"`
		LastName  string `validate:"required" json:"last_name"`
		Email     string `validate:"required,email" json:"email"`
		Password  string `validate:"required" json:"password"`
	}

	User struct {
		UserID    string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
		FirstName string    `json:"first_name" bson:"first_name"`
		LastName  string    `json:"last_name" bson:"last_name"`
		Email     string    `json:"email" bson:"email"`
		Password  string    `json:"password,omitempty" bson:"password,omitempty"`
		Locked    bool      `json:"locked,omitempty" bson:"locked,omitempty"`
		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
		AvatarURL string    `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
		Roles     []string  `json:"roles,omitempty" bson:"roles,omitempty"`
	}

	// UserInfo struct {
	// 	Email     string     `json:"email,omitempty"`
	// 	FirstName string     `json:"first_name,omitempty"`
	// 	LastName  string     `json:"last_name,omitempty"`
	// 	UserID    string     `json:"user_id,omitempty"`
	// 	AvatarURL string     `json:"avatar_url,omitempty"`
	// 	CreatedAt *time.Time `json:"created_at,omitempty"`
	// 	UpdateAt  *time.Time `json:"updated_at,omitempty"`
	// }
)

func (user *User) Strip() *User {
	stripedUser := User(*user)
	stripedUser.Password = ""
	return &stripedUser
}
