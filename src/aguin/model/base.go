package model

import (
	"aguin/config"
	"gopkg.in/mgo.v2"
)

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
