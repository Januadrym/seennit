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
		ID        string    `bson:"_id"`
		FirstName string    `bson:"first_name"`
		LastName  string    `bson:"last_name"`
		Email     string    `bson:"email"`
		Password  string    `bson:"password"`
		Locked    bool      `bson:"locked"`
		CreatedAt time.Time `bson:"created_at"`
		UpdatedAt time.Time `bson:"updated_at"`
	}
)

func (user *User) Strip() *User {
	stripedUser := User(*user)
	stripedUser.Password = ""
	return &stripedUser
}
