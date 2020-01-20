package user

import (
	"context"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
)

type (
	MongoDBRepository struct {
		session *mgo.Session
		Users   map[string]User
	}
)

func NewMongoDBRepo(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session: session,
	}
}

func (r *MongoDBRepository) Insert(ctx context.Context, user User) error {
	s := r.session.Clone()
	defer s.Clone()
	if err := s.DB("").C("users").Insert(user); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) Create(ctx context.Context, user User) (string, error) {
	id := uuid.New().String()
	r.Users[id] = user
	return id, nil
}

func (r *MongoDBRepository) Home(ctx context.Context, user User) (string, error) {
	name := "this is Q"
	r.Users[name] = user
	return name, nil
}

func (r *MongoDBRepository) FindUserByMail(ctx context.Context, email string) (*User, error) {
	s := r.session.Clone()
	defer s.Clone()
	var usr User
	if err := s.DB("").C("user").Find(bson.M{
		"email": email,
	}).One(&usr); err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &usr, nil
}
