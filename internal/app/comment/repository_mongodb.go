package comment

import (
	"context"

	"github.com/Januadrym/seennit/internal/app/types"
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
