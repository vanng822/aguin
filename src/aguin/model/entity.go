package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Entity struct {
	Id        bson.ObjectId          `json:"-" bson:"_id,omitempty"`
	Name      string                 `json:"entity"`
	AppId     bson.ObjectId          `json:"-"`
	CreatedAt time.Time              `json:"createdAt"`
	Data      map[string]interface{} `json:"data"`
}

func (en *Entity) Save(session *mgo.Session) error {
	return EntityCollection(session).Insert(en)
}
