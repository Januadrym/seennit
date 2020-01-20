package user

import (
	"context"
	"github.com/google/uuid"
)

type InMemoryRepo struct {
	Users map[string]User
}

func (r *InMemoryRepo) Insert(ctx context.Context, user User) error {
	return nil
}

func (r *InMemoryRepo) Create(ctx context.Context, user User) (string, error) {
	id := uuid.New().String()
	r.Users[id] = user
	return id, nil
}

func (r *InMemoryRepo) Home(ctx context.Context, user User) (string, error) {
	name := "this is Q"
	r.Users[name] = user
	return name, nil
}

func (r *InMemoryRepo) FindUserByName(ctx context.Context, name string) string {
	if name == "Q" {
		return "same"
	}
	return "OK"
}
