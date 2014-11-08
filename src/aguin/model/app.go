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

/*
func (m *Application) serialize() map[string]interface{} {
	return map[string]interface{}{
			"name": m.Name,
			"appid": }
}*/
/*
func init() {
	c := Collection("app")
	c.EnsureIndex(mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	})
	c.EnsureIndex(mgo.Index{
		Key:        []string{"Userid"},
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	})
}*/


func (m *Application) Save(session *mgo.Session) error {
	return AppCollection(session).Insert(m)
}

func Get(fromDate, toDate time.Time) interface{} {
	return nil
}
