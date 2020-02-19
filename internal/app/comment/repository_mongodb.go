package comment

import (
	"context"

	"github.com/Januadrym/seennit/internal/app/types"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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
	return s.DB("").C("comments")
}

func (r *MongoDBRepository) FindAll(ctx context.Context) ([]*types.Comment, error) {
	s := r.sessions.Clone()
	defer s.Close()
	var list []*types.Comment
	if err := r.collection(s).Find(nil).All(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *MongoDBRepository) Create(ctx context.Context, comment *types.Comment) error {
	s := r.sessions.Clone()
	defer s.Close()
	if err := r.collection(s).Insert(&comment); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) FindCommentPost(ctx context.Context, idPost string) ([]*types.Comment, error) {
	s := r.sessions.Clone()
	defer s.Close()
	var list []*types.Comment
	if err := r.collection(s).Find(bson.M{"post_id": idPost}).All(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *MongoDBRepository) Update(ctx context.Context, id string, c string) error {
	s := r.sessions.Clone()
	defer s.Close()
	return r.collection(s).Update(bson.M{"id": id}, bson.M{
		"$set": bson.M{
			"content": c, // content of a comment
		},
	})
}

func (r *MongoDBRepository) DeleteByID(ctx context.Context, id string) error {
	s := r.sessions.Clone()
	defer s.Close()
	return r.collection(s).Remove(bson.M{"id": id})
}
