package community

import (
	"context"

	"github.com/Januadrym/seennit/internal/app/types"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type (
	MongoDBRepository struct {
		session *mgo.Session
	}
)

func NewMongoDBRepo(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session: session,
	}
}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("communities")
}

func (r *MongoDBRepository) FindCommunityByID(ctx context.Context, cID string) (*types.Community, error) {
	s := r.session.Clone()
	defer s.Clone()
	var com *types.Community
	if err := r.collection(s).Find(bson.M{
		"communityID": cID,
	}).One(&com); err != nil {
		return nil, err
	}
	return com, nil
}

func (r *MongoDBRepository) FindAll(context.Context) ([]*types.Community, error) {
	s := r.session.Clone()
	defer s.Close()
	var coms []*types.Community
	if err := r.collection(s).Find(nil).All(&coms); err != nil {
		return nil, err
	}
	return coms, nil
}
