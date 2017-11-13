package run_test

import (
	"os"
	"testing"

	"github.com/user/numb/run"
)

func TestRun(t *testing.T) {
	run.Train("mkdir temp", nil)
	if _, err := os.Stat("temp"); os.IsNotExist(err) {
		t.Error("dir not really created!")
	}
	run.Test("rmdir temp", map[string]interface{}{"silent": true})
	if _, err := os.Stat("temp"); os.IsExist(err) {
		t.Error("dir not really removed!")
	}
}
