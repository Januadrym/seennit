package post

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
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

func (r *MongoDBRepository) Create(ctx context.Context, req *Post) error {
	s := r.sessions.Clone()
	defer s.Close()
	req.UpdatedAt = req.CreatedAt

	if err := r.collection(s).Insert(req); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) GetAll(ctx context.Context, listID []string) ([]*Post, error) {
	s := r.sessions.Clone()
	defer s.Close()
	var listPost []*Post
	if err := r.collection(s).Find(bson.M{"id": bson.M{
		"$in": listID,
	}}).All(&listPost); err != nil {
		logrus.Errorf("failed to get posts, err : %v", err)
		return nil, err
	}
	return listPost, nil
}
