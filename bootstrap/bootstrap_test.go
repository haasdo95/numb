package bootstrap_test

import (
	"testing"

	"github.com/user/numb/bootstrap"

	"github.com/user/numb/run"
)

func TestInit(t *testing.T) {
	run.Train("numb deinit", nil)
	errdir, errjson := bootstrap.Init()
	if errdir != nil || errjson != nil {
		t.Error("An Error Occurred When Init")
	}
	errdir, errjson = bootstrap.Init()
	if errdir == nil || errjson == nil {
		t.Error("Error Not Caught When Init")
	}
	errdir, errjson = bootstrap.Deinit()
	if errdir != nil || errjson != nil {
		t.Error("An Error Occurred When Deinit")
	}
	errdir, errjson = bootstrap.Deinit()
	if errdir == nil || errjson == nil {
		t.Error("Error Not Caught When Deinit")
	}
}

func TestDeinit(t *testing.T) {

}
