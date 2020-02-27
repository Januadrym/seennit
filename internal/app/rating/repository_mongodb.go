package rating

import (
	"context"
	"time"

	"github.com/Januadrym/seennit/internal/app/types"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
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

func (r *MongoDBRepository) collections(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("rating")
}

// Upsert update and insert :\
func (r *MongoDBRepository) Upsert(ctx context.Context, rating *types.Rating) (bool, error) {
	s := r.sessions.Clone()
	idUpdate := uuid.New().String()
	isNew := false
	defer s.Close()
	selector := bson.M{
		"target_type":   rating.TargetType,
		"target_id":     rating.TargetID,
		"created_by_id": rating.CreatedByID,
	}
	updates := bson.M{
		"$set": bson.M{
			"type":            rating.Type,
			"created_by_id":   rating.CreatedByID,
			"created_by_name": rating.CreatedByName,
			"created_at":      time.Now(),
		},
		"$setOnInsert": bson.M{
			"id":          idUpdate,
			"target_type": rating.TargetType,
			"target_id":   rating.TargetID,
		},
	}
	inf, err := r.collections(s).Upsert(selector, updates)
	logrus.Infof("Status: MATCHED: %v, REMOVED: %v, ID: %v, UPDATED: %v ", inf.Matched, inf.Removed, inf.UpsertedId, inf.Updated)
	if err != nil {
		logrus.Errorf("failed to create/update rating info, err: %v", err)
		return false, err
	}
	if inf.UpsertedId != nil {
		isNew = true
	}
	return isNew, nil
}
