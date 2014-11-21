package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id    bson.ObjectId `json:"id" bson:"_id,omitempty"`
	FbId  int64         `json:"-" bson:"fbid"`
	Name  string        `json:"name" bson:"name"`
	Email string        `json:"-" bson:"email"`
}

func (u *User) Save(session *mgo.Session) error {
	return UserCollection(session).Insert(u)
}
