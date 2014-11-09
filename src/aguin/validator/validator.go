package validator

import (
	"log"
	"regexp"
	"time"
)

var allowedChars = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
var validObjectId = regexp.MustCompile(`^[a-f0-9]{24}$`)

func ValidateObjectId(id string) bool {
	return validObjectId.MatchString(id)
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

func ValidateEntity(message map[string]interface{}) (string, map[string]interface{}, bool) {
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

	entity = ValidateEntityName(entity)

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
			newData[k] = newArr
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
		return entity, map[string]interface{}{}, false
	}
	return entity, newData, len(data) > 0
}
