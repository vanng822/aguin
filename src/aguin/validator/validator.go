package validator

import (
	"regexp"
	"time"
	"aguin/utils"
)

var allowedChars = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
var validAPIKey = regexp.MustCompile(`^[a-f0-9]{24}$`)
const dateInputFormat = "2006-01-02"

func ValidAPIKey(id string) bool {
	return validAPIKey.MatchString(id)
}

func ValidateEntityName(name string) string {
	if allowedChars.MatchString(name) {
		return name
	}
	return ""
}

type SearchSchema struct {
	Validated bool
	Entity    string
	StartDate time.Time
	EndDate   time.Time
}


func ValidateSearch(message map[string]interface{}) SearchSchema {
	log := utils.GetLogger("aguin")
	entity, _ := message["entity"].(string)
	entity = ValidateEntityName(entity)
	validated := entity != ""
	schema := SearchSchema{}
	schema.Validated = validated
	schema.Entity = entity
	// if dates invalid set to latest 30 days
	var start, end time.Time
	var err error
	
	startDate, ok := message["startDate"].(string)
	log.Debug("startDate:%v,%v", startDate, ok)
	
	if ok {
		start, err = time.Parse(dateInputFormat, startDate)
		log.Debug("start:%v,%v", start, err)
	}
	endDate, ok := message["endDate"].(string)
	log.Debug("endDate:%v,%v", endDate, ok)
	if ok {
		end, err = time.Parse(dateInputFormat, endDate)
		log.Debug("end:%v,%v", end, err)
	}
	
	if !start.IsZero() && !end.IsZero() {
		if start.After(end) {
			end = time.Now()
			start = end.AddDate(0, 0, -30)
		} else if end.AddDate(0, 0, -90).After(start) {
			start = end.AddDate(0, 0, -90)
		}
		
	} else {
		end = time.Now()
		start = schema.EndDate.AddDate(0, 0, -30)
	}
	
	// set start date time at 00:00:00 and end date to 23:59:59
	schema.StartDate = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	schema.EndDate = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())
	
	return schema
}

type EntitySchema struct {
	Validated bool
	Entity    string
	Data      map[string]interface{}
}

func ValidateEntity(message map[string]interface{}) EntitySchema {
	/*Allow type:
	int
	float
	date: 2006-01-02 03:04:01
	boolean
	array[int, float]
	*/
	errCounter := 0
	data, _ := message["data"].(map[string]interface{})
	entity, _ := message["entity"].(string)
	newData := map[string]interface{}{}

	result := EntitySchema{}
	result.Entity = ValidateEntityName(entity)
	for k, v := range data {
		if !allowedChars.MatchString(k) {
			errCounter += 1
			continue
		}
		switch vv := v.(type) {
		case int:
			newData[k] = vv
		case float64:
			newData[k] = vv
		case bool:
			newData[k] = vv
		case []interface{}:
			newArr := make([]interface{}, 0)
			for _, u := range vv {
				switch u.(type) {
				case int:
					newArr = append(newArr, u)
				case float64:
					newArr = append(newArr, u)
				default:
					errCounter += 1
				}
			}
			if len(newArr) > 0 {
				newData[k] = newArr
			}
		default:
			errCounter += 1
		}
	}
	if errCounter > 0 {
		result.Validated = false
		result.Data = map[string]interface{}{}
	} else {
		result.Validated = len(newData) > 0
		result.Data = newData
	}
	return result
}
