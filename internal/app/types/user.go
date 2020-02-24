package types

import (
	"fmt"
	"strings"
	"time"
)

type (
	RegisterRequest struct {
		FirstName string `validate:"required" json:"first_name"`
		LastName  string `validate:"required" json:"last_name"`
		Email     string `validate:"required,email" json:"email"`
		Password  string `validate:"required,gt=3" json:"password"`
	}

	User struct {
		UserID    string `json:"id,omitempty" bson:"id,omitempty"`
		FirstName string `validate:"required" json:"first_name" bson:"first_name"`
		LastName  string `validate:"required" json:"last_name" bson:"last_name"`
		AvatarURL string `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`

		Email    string `validate:"required" json:"email" bson:"email"`
		Password string `json:"password,omitempty" bson:"password,omitempty"`

		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`

		Communities []string `json:"communities,omitempty" bson:"communities"`

		KarmaPoints  uint32   `json:"karma_points,omitempty" bson:"karma_points,omitempty"`
		Achievements []string `json:"achievements,omitempty" bson:"achievements,omitempty"`
	}
)

func (user *User) Strip() *User {
	stripedUser := User(*user)
	stripedUser.Password = ""
	return &stripedUser
}

func (user *User) FullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

// GetName should always be used to get the user name
func (user User) GetName() string {
	if user.FullName() != "" {
		return user.FullName()
	}
	return strings.Split(user.Email, "@")[0]
}
