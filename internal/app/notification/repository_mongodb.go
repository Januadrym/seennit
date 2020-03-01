package notification

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
	return s.DB("").C("notification")
}

func (r *MongoDBRepository) Create(ctx context.Context, noti *types.PushNotification) error {
	s := r.sessions.Clone()
	defer s.Close()
	if err := r.collection(s).Insert(noti); err != nil {
		return err
	}
	return nil
}
