package utils_test

import (
	"io/ioutil"
	"testing"

	"github.com/user/numb/utils"
)

func TestCon2Abs(t *testing.T) {
	raw, err := ioutil.ReadFile("./example.proto")
	utils.Check(err)
	con := string(raw)
	_, abs := utils.Concrete2Abstract(con)
	t.Logf("ABS: %s", abs)
}
