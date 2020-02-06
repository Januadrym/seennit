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
	defer s.Close()
	var com *types.Community
	if err := r.collection(s).Find(bson.M{
		"community_ID": cID,
	}).One(&com); err != nil {
		return nil, err
	}
	return com, nil
}

func (r *MongoDBRepository) FindCommunityByName(ctx context.Context, cName string) (*types.Community, error) {
	s := r.session.Clone()
	defer s.Close()
	var com *types.Community
	if err := r.collection(s).Find(bson.M{
		"community_name": cName,
	}).One(&com); err != nil {
		return nil, err
	}
	return com, nil
}

func (r *MongoDBRepository) FindAllCom(context.Context) ([]*types.Community, error) {
	s := r.session.Clone()
	defer s.Close()
	var coms []*types.Community
	if err := r.collection(s).Find(nil).All(&coms); err != nil {
		return nil, err
	}
	return coms, nil
}

func (r *MongoDBRepository) Create(ctx context.Context, com *types.Community) error {
	s := r.session.Clone()
	defer s.Close()
	com.UpdatedAt = com.CreatedAt

	if err := r.collection(s).Insert(com); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) Delete(ctx context.Context) error {
	s := r.session.Clone()
	defer s.Close()
	r.collection(s).RemoveAll(nil)
	return nil
}
