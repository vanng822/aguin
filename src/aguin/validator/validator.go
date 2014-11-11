package validator

import (
	"log"
	"regexp"
	"time"
)

var allowedChars = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
var validAPIKey = regexp.MustCompile(`^[a-f0-9]{24}$`)

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
}

func ValidateSearch(message map[string]interface{}) SearchSchema {
	entity, _ := message["entity"].(string)
	entity = ValidateEntityName(entity)

	return SearchSchema{entity != "", entity}
}

type EntitySchema struct {
	Validated bool
	Entity string
	Data map[string]interface{}
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
		case string:
			vvv, err := time.Parse("2006-01-02 03:04:01", vv)
			if err != nil {
				log.Println(err)
				errCounter += 1
			} else {
				newData[k] = vvv
			}
		default:
			errCounter += 1
		}
	}
	if errCounter > 0 {
		result.Validated = false
	} else {
		result.Validated = len(newData) > 0
	}
	return result
}
