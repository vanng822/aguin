package model

import (
	"aguin/utils"
	"gopkg.in/mgo.v2"
)

func EnsureIndex(concurrent bool) {
	log := utils.GetLogger("system")
	if concurrent {
		log.Info("EnsureIndex running in background")
		go ensureUserIndex()
		go ensureAppIndex()
		go ensureEntityIndex()
	} else {
		log.Info("EnsureIndex running")
		ensureUserIndex()
		ensureAppIndex()
		ensureEntityIndex()
	}
}

func countIndexes(c *mgo.Collection) int {
	log := utils.GetLogger("system")
	exist_indexes, err := c.Indexes()
	if err != nil {
		log.Warning("Error when fetching indexes info for collection %s: %v", c.Name, err)
		return 0
	}

	return len(exist_indexes)
}

func ensureIndexes(indexes []mgo.Index, c *mgo.Collection) {
	log := utils.GetLogger("system")
	// naive check, minus one for primary _id
	index_count := countIndexes(c) - 1
	if len(indexes) != index_count {
		log.Info("Will run ensureIndex for collection: %s, new: %d, existing: %d", c.Name, len(indexes), index_count)
		for _, index := range indexes {
			c.EnsureIndex(index)
		}
	} else {
		log.Info("Will not run ensureIndex for collection: %s", c.Name)
	}
}

func ensureEntityIndex() {
	session := Session()
	defer session.Close()

	indexes := []mgo.Index{
		mgo.Index{
			Key:        []string{"name"},
			Background: true},
		mgo.Index{
			Key:        []string{"appid"},
			Background: true},
		mgo.Index{
			Key:        []string{"created_at"},
			Background: true}}

	ensureIndexes(indexes, EntityCollection(session))
}

func ensureAppIndex() {
	session := Session()
	defer session.Close()
	indexes := []mgo.Index{
		mgo.Index{
			Key:        []string{"name", "userid"},
			Unique:     true,
			Sparse:     true,
			Background: true}}
	ensureIndexes(indexes, AppCollection(session))
}

func ensureUserIndex() {
	session := Session()
	defer session.Close()

	indexes := []mgo.Index{
		mgo.Index{
			Key:        []string{"email"},
			Unique:     true,
			Sparse:     true,
			Background: true}}

	ensureIndexes(indexes, UserCollection(session))
}
