package api

import (
	"sync"

	"github.com/globalsign/mgo"
)

var (
	session     *mgo.Session
	sessionOnce sync.Once
)

// func dialMongo() (*mgo.Session, error) {
// 	repoConf, _ := mongodb.Load()
// 	var err error
// 	sessionOnce.Do(func() {
// 		session, err = mongodb.Dial(repoConf)
// 	})
// }
