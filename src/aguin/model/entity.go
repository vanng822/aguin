package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Entity struct {
	Id        bson.ObjectId          `json:"id" bson:"_id,omitempty"`
	Name      string                 `json:"entity" bson:"name"`
	AppId     bson.ObjectId          `json:"-" bson:"appid"`
	CreatedAt time.Time              `json:"createdAt" bson:"createdat"`
	Data      map[string]interface{} `json:"data" bson:"data"`
}

func (e *Entity) Save(session *mgo.Session) error {
	return EntityCollection(session).Insert(e)
}