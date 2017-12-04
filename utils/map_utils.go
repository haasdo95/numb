package utils

import (
	"sort"
	"encoding/json"
)

// MapKeyValue returns (keys, values) of a map
// keys will always be returned in increasing order
func MapKeyValue(m map[string]interface{}) ([]string, []interface{}) {
	keys := make([]string, len(m))
	values := make([]interface{}, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for j := 0; j < i; j++ {
		values[j] = m[keys[j]]
	}
	return keys, values
}

func Str2Obj(jsonStr string) (map[string]interface{}, error) {
	var obj map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &obj)
	return obj, err
}