package utils

import (
	"encoding/json"
)

// MapKeyValue returns (keys, values) of a map
func MapKeyValue(m map[string]interface{}) ([]string, []interface{}) {
	keys := make([]string, len(m))
	values := make([]interface{}, len(m))
	i := 0
	for k, v := range m {
		keys[i] = k
		values[i] = v
		i++
	}
	return keys, values
}

func Str2Obj(jsonStr string) (map[string]interface{}, error) {
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	return obj, err
}