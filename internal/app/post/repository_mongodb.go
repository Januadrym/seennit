package post

import (
	"context"
	"time"

	"github.com/Januadrym/seennit/internal/app/types"
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

func (r *MongoDBRepository) Create(ctx context.Context, req *types.Post) error {
	s := r.sessions.Clone()
	defer s.Close()
	req.UpdatedAt = req.CreatedAt

	if err := r.collection(s).Insert(req); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) FindByID(ctx context.Context, id string) (*types.Post, error) {
	s := r.sessions.Clone()
	defer s.Close()
	var p *types.Post
	if err := r.collection(s).Find(bson.M{"id": id}).One(&p); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *MongoDBRepository) GetAll(ctx context.Context, listID []string) ([]*types.Post, error) {
	s := r.sessions.Clone()
	defer s.Close()
	var listPost []*types.Post
	if err := r.collection(s).Find(bson.M{"id": bson.M{
		"$in": listID,
	}}).All(&listPost); err != nil {
		logrus.Errorf("failed to get posts, err : %v", err)
		return nil, err
	}
	return listPost, nil
}

func (r *MongoDBRepository) UpdatePost(ctx context.Context, id string, p *types.PostUpdateRequest) error {
	s := r.sessions.Clone()
	defer s.Close()
	p.UpdatedAt = time.Now()
	logrus.Info("post mongo ne:", p)
	if err := r.collection(s).Update(bson.D{{Name: "id", Value: id}}, bson.M{"$set": p}); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) ChangeStatus(ctx context.Context, id string, status types.Status) error {
	s := r.sessions.Clone()
	defer s.Close()
	if err := r.collection(s).Update(bson.M{"id": id}, bson.M{"$set": bson.M{"status": status}}); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) GetEntire(ctx context.Context) ([]*types.Post, error) {
	s := r.sessions.Clone()
	defer s.Close()
	var list []*types.Post
	if err := r.collection(s).Find(nil).All(&list); err != nil {
		return nil, err
	}
	return list, nil
}
