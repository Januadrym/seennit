package community

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
	return s.DB("").C("communities")
}

func (r *MongoDBRepository) FindCommunityByID(ctx context.Context, cID string) (*types.Community, error) {
	s := r.session.Clone()
	defer s.Close()
	var com *types.Community
	if err := r.collection(s).Find(bson.M{
		"id": cID,
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
		"name": cName,
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

// redo later
func (r *MongoDBRepository) DeleteByID(ctx context.Context, id string) error {
	s := r.session.Clone()
	defer s.Close()
	return r.collection(s).Remove(bson.M{"id": id})
}

func (r *MongoDBRepository) EnrollUser(ctx context.Context, idUser string, idCom string) error {
	s := r.session.Clone()
	defer s.Close()

	if err := r.collection(s).Update(bson.M{"id": idCom}, bson.M{
		"$addToSet": bson.M{
			"users": idUser,
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

func (r *MongoDBRepository) UpdateInfo(ctx context.Context, idCom string, comm *types.Community) error {
	s := r.session.Clone()
	defer s.Close()
	return r.collection(s).Update(bson.M{"id": idCom}, bson.M{
		"$set": bson.M{
			"name":        comm.Name,
			"description": comm.Description,
			"banner_URL":  comm.BannerURL,
			"updated_at":  time.Now(),
		},
	})
}

func (r *MongoDBRepository) AddPost(ctx context.Context, idPost string, idCom string) error {
	s := r.session.Clone()
	defer s.Close()
	if err := r.collection(s).Update(bson.M{"id": idCom}, bson.M{
		"$addToSet": bson.M{
			"posts": idPost,
		},
	}); err != nil {
		logrus.Errorf("failed to added post, err : %v", err)
		return err
	}
	return nil
}

func (r *MongoDBRepository) GetAllPost(ctx context.Context, idCom string) ([]string, error) {
	s := r.session.Clone()
	defer s.Close()
	var com *types.Community
	if err := r.collection(s).Find(bson.M{"id": idCom}).One(&com); err != nil {
		return nil, err
	}
	return com.Posts, nil
}

func (r *MongoDBRepository) CheckContainPost(ctx context.Context, comName, idPost string) (bool, error) {
	s := r.session.Clone()
	defer s.Close()
	var com *types.Community
	if err := r.collection(s).Find(bson.M{"name": comName, "posts": idPost}).One(&com); err != nil {
		return false, err
	}
	return true, nil
}
