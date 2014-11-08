package model

import (
	"aguin/config"
	"gopkg.in/mgo.v2"
	"time"
)

func Test() (mgo.BuildInfo, error) {
	conf := config.AppConf()
	session, err := mgo.Dial(conf.Mongodb)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	return session.BuildInfo()
}

func Session() *mgo.Session {
	conf := config.AppConf()
	session, err := mgo.Dial(conf.Mongodb)
	if err != nil {
		panic(err)
	}
	return session
}

func UserCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("").C("user")
}

func AppCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("").C("app")
}

func EntityCollection(session *mgo.Session) *mgo.Collection {
	return session.DB("").C("entity")
}


type IModel interface {
	Save(data interface{}) error
	Get(fromDate, toDate time.Time) interface{}
	Update(date time.Time, data map[string]interface{}) error
}
