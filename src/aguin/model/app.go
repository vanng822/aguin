package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Application struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserId      bson.ObjectId `json:"userid" bson:"userid"`
	Name        string        `json:"name" bson:"name"`
	Secret      string        `json:"-" bson:"secret"`
	Description string        `json:"description" bson:"description"`
}

func (a *Application) Save(session *mgo.Session) error {
	return AppCollection(session).Insert(a)
}
