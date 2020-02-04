package user

import (
	"context"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"

	"github.com/Januadrym/seennit/internal/app/types"
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

func (r *MongoDBRepository) Insert(ctx context.Context, user *types.User) error {
	logrus.Errorf("insert user: %v", user)
	s := r.session.Clone()
	defer s.Clone()
	if err := r.collection(s).Insert(user); err != nil {
		logrus.Errorf("insert user: %v", user)
		return err
	}
	return nil
}

func (r *MongoDBRepository) FindUserByMail(ctx context.Context, email string) (*types.User, error) {
	s := r.session.Clone()
	defer s.Clone()
	var usr *types.User
	if err := r.collection(s).Find(bson.M{
		"email": email,
	}).One(&usr); err != nil {
		return nil, err
	}
	return usr, nil
}

func (r *MongoDBRepository) FindAll(context.Context) ([]*types.User, error) {
	s := r.session.Clone()
	defer s.Close()
	var users []*types.User
	if err := r.collection(s).Find(nil).All(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *MongoDBRepository) Delete(ctx context.Context) error {
	s := r.session.Clone()
	defer s.Close()
	r.collection(s).RemoveAll(nil)
	return nil
}

func (r *MongoDBRepository) UpdateInfo(ctx context.Context, userID string, user *types.User) error {
	s := r.session.Clone()
	defer s.Clone()
	return r.collection(s).Update(bson.M{"user_id": userID}, bson.M{
		"$set": bson.M{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"updated_at": time.Now(),
		},
	})
}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("users")
}
