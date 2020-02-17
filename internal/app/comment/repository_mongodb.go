package comment

import "github.com/globalsign/mgo"

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
