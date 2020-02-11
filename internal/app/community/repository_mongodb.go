package community

import (
	"context"

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

func (r *MongoDBRepository) FindCommunityByID(ctx context.Context, cID string) (*Community, error) {
	s := r.session.Clone()
	defer s.Close()
	var com *Community
	if err := r.collection(s).Find(bson.M{
		"ID": cID,
	}).One(&com); err != nil {
		return nil, err
	}
	return com, nil
}

func (r *MongoDBRepository) FindCommunityByName(ctx context.Context, cName string) (*Community, error) {
	s := r.session.Clone()
	defer s.Close()
	var com *Community
	if err := r.collection(s).Find(bson.M{
		"name": cName,
	}).One(&com); err != nil {
		return nil, err
	}
	return com, nil
}

func (r *MongoDBRepository) FindAllCom(context.Context) ([]*Community, error) {
	s := r.session.Clone()
	defer s.Close()
	var coms []*Community
	if err := r.collection(s).Find(nil).All(&coms); err != nil {
		return nil, err
	}
	return coms, nil
}

func (r *MongoDBRepository) Create(ctx context.Context, com *Community) error {
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

// redo later
func (r *MongoDBRepository) DeleteByID(ctx context.Context, id string) error {
	s := r.session.Clone()
	defer s.Close()
	return r.collection(s).Remove(bson.M{"ID": id})
}

func (r *MongoDBRepository) EnrollUser(ctx context.Context, idUser string, idCom string) error {
	s := r.session.Clone()
	defer s.Close()
	return r.collection(s).Update(bson.M{"ID": idCom}, bson.M{
		"$addToSet": bson.M{
			"user": idUser,
		},
	})
}
