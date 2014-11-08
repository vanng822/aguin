package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Application struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Userid      bson.ObjectId
	Name        string
	Secret      string 
	Description string
}

func (m *Application) Save(session *mgo.Session) error {
	return AppCollection(session).Insert(m)
}
