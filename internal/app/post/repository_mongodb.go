package post

import (
	"context"

	"github.com/globalsign/mgo"
)

type (
	MongoDBRepository struct {
		sessions *mgo.Session
	}
)

func NewMongoDBRepo(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		sessions: session,
	}
}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("posts")
}

func (r *MongoDBRepository) Create(ctx context.Context, post *Post) error {
	s := r.sessions.Clone()
	defer s.Close()
	post.UpdatedAt = post.CreatedAt

	if err := r.collection(s).Insert(post); err != nil {
		return err
	}
	return nil
}
