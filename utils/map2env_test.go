package utils_test

import (
	"fmt"
	"testing"

	"../utils"
)

func TestMap2Env(t *testing.T) {
	m1 := map[string]string{}
	m2 := map[string]string{
		"MODE":   "TRAIN",
		"RECORD": "FALSE",
	}

	ss1 := utils.Map2env(m1)
	ss2 := utils.Map2env(m2)
	fmt.Println((ss1))
	fmt.Println((ss2))
	if len(ss1) != 0 {
		t.Error("Expecting Empty Slice")
	}
	if len(ss2) != 2 {
		t.Error("Expecting String Slice of len 2")
	}

}
