package utils

import (
	"fmt"
)

// Map2env converts a map into a slice of <env>=<val> strings
func Map2env(m map[string]string) []string {
	env := make([]string, 0)
	for k, v := range m {
		s := fmt.Sprintf("%s=%s", k, v)
		env = append(env, s)
	}
	return env
}
