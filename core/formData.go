package core

import (
	"encoding/json"
	"strings"
)

func formatFormFieldValue(value string) (interface{}, error) {
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		value = strings.TrimPrefix(value, "[")
		value = strings.TrimSuffix(value, "]")
		isMap := strings.Contains(value, "{")
		splittedValue := strings.Split(value, ",")
		resultSlice := make([]interface{}, len(splittedValue))
		resultMaps := make([]map[string]interface{}, len(splittedValue))

		for i, item := range splittedValue {
			if strings.HasPrefix(item, "{") && strings.HasSuffix(item, "}") {
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(item), &data); err != nil {
					return data, err
				}
				resultMaps[i] = data
			} else {
				resultSlice[i] = item
			}
		}

		if isMap {
			return resultMaps, nil
		}
		return resultSlice, nil
	}

	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(value), &value); err != nil {
			return data, err
		}
		return data, nil
	}

	return value, nil
}
