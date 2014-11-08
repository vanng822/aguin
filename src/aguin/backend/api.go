package backend

import (
	"time"
)

type Backend interface {
	Save(data interface{}) bool
	Get(fromDate, toDate time.Time) interface{}
	Update(date time.Time, data map[string]interface{}) bool
}

func GetBackend(factory string) Backend {
	switch factory {
	case "mongodb":
		return MongodbBackend{}
	default:
		return MongodbBackend{}
	}
}
