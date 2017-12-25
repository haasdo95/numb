package bootstrap_test

import (
	"testing"

	"github.com/haasdo95/numb/bootstrap"

	"os/exec"
)

func TestInit(t *testing.T) {
	exec.Command("numb", "deinit").Run()
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
