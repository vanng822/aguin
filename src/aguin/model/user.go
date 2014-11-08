package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Fbid int64
	Name string
	Email string
}


func (u *User) Save(session *mgo.Session) error {
	return UserCollection(session).Insert(u)
}
