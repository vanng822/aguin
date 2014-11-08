package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Application struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Userid      bson.ObjectId
	Name        string
	Secret      string `json:"-"`
	Description string
}

func (m *Application) Save(session *mgo.Session) error {
	return AppCollection(session).Insert(m)
}
