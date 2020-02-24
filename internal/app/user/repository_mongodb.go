package user

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
		session *mgo.Session
	}
)

func NewMongoDBRepo(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session: session,
	}
}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("users")
}

func (r *MongoDBRepository) Create(ctx context.Context, user *types.User) error {
	s := r.session.Clone()
	defer s.Close()
	user.UpdatedAt = user.CreatedAt

	if err := r.collection(s).Insert(user); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) FindUserByMail(ctx context.Context, email string) (*types.User, error) {
	s := r.session.Clone()
	defer s.Close()
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

func (r *MongoDBRepository) Delete(ctx context.Context, id string) error {
	s := r.session.Clone()
	defer s.Close()
	return r.collection(s).Remove(bson.M{"id": id})
}

func (r *MongoDBRepository) UpdateInfo(ctx context.Context, userID string, user *types.User) error {
	s := r.session.Clone()
	defer s.Close()
	return r.collection(s).Update(bson.M{"id": userID}, bson.M{
		"$set": bson.M{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"avatar_url": user.AvatarURL,
			"updated_at": time.Now(),
		},
	})
}

func (r *MongoDBRepository) EnrollUser(ctx context.Context, idUser string, idCom string) error {
	s := r.session.Clone()
	defer s.Close()

	if err := r.collection(s).Update(bson.M{"id": idUser}, bson.M{
		"$addToSet": bson.M{
			"communities": idCom,
		},
	}); err != nil {
		logrus.Errorf("failed to enroll user, err: %v", err)
		return err
	}
	return nil
}

func (r *MongoDBRepository) CheckUserEnrolled(ctx context.Context, idUser string, idCom string) (string, error) {
	s := r.session.Clone()
	defer s.Close()
	var com *types.Community
	if err := r.collection(s).Find(bson.M{"id": idCom, "users": idUser}).One(&com); err != nil {
		return "", err
	}
	return idUser, nil
}
