package jsonHelper

import (
	"encoding/json"
	"log"
)

func ParseJSON(p interface{}) []byte {
	payload, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	return payload
}

func GetJSONFromStruct(s interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}

	bytes, err := json.Marshal(s)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
