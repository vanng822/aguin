package utils

import (
	"encoding/json"
)

func Bytes2json(data *[]byte) (map[string]interface{}, error) {
	var jsonData interface{}

	err := json.Unmarshal(*data, &jsonData)

	if err != nil {
		return nil, err
	}

	return jsonData.(map[string]interface{}), nil
}

func Json2bytes(data interface{}) (*[]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &b, nil
}