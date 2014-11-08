package backend

import (
	"gopkg.in/mgo.v2"
	"aguin/config"
	"aguin/utils"
	"time"
)

type MongodbBackend struct {
	
}

func (m MongodbBackend) Save(entity interface{}) bool {
	log := utils.GetLogger("mongodb")
	log.Print("Backend Save")
	log.Print(entity)
	return true
}

func (m MongodbBackend) Get(fromDate, toDate time.Time) interface{} {
	return nil
}

func (m MongodbBackend) Update(date time.Time, data map[string]interface{}) bool {
	return false
}