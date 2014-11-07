package utils

import (
	"log"
	"time"
	//"reflect"
)

func Validate(data map[string]interface{}) (map[string]interface{}, bool) {
	/*Allow type:
		int
		float
		date: 2006-01-02 03:04:01
		boolean
		array[int, float]
	*/
	errCounter := 0
	newData := map[string]interface{}{}
	
	for k, v := range data {
		switch vv := v.(type) {
		case int:
			newData[k] = vv
		case float64:
			newData[k] = vv
		case bool:
			newData[k] = vv
		case []interface{}:
			newArr := make([]interface{},0)
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
		return map[string]interface{}{}, false
	}
	return newData, true
}
