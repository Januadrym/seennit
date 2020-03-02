package notification

import (
	"context"

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
	return s.DB("").C("notification")
}

func (r *MongoDBRepository) Create(ctx context.Context, noti *types.PushNotification) error {
	s := r.sessions.Clone()
	defer s.Close()
	if err := r.collection(s).Insert(noti); err != nil {
		logrus.Errorf("failed to create notificaion in DB, err: %v", err)
		return err
	}
	return nil
}

func (r *MongoDBRepository) LoadNoti(ctx context.Context, userID string) ([]*types.PushNotification, error) {
	s := r.sessions.Clone()
	defer s.Close()
	var notis []*types.PushNotification
	if err := r.collection(s).Find(bson.M{"user_id": userID}).All(&notis); err != nil {
		logrus.Errorf("failed to find notification in DB, err: %v", err)
		return nil, err
	}
	return notis, nil
}
