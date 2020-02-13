package api

import (
	"sync"

	"github.com/Januadrym/seennit/internal/pkg/db/mongodb"
	"github.com/globalsign/mgo"
)

var (
	session     *mgo.Session
	sessionOnce sync.Once
)

func dialMongo() (*mgo.Session, error) {
	repoConf, _ := mongodb.Load()
	var err error
	sessionOnce.Do(func() {
		session, err = mongodb.Dial(repoConf)
	})
	if err != nil {
		return nil, err
	}
	s := session.Clone()
	return s, nil
}
