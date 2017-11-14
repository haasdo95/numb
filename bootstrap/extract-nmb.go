package bootstrap

import (
	"encoding/json"
	"io/ioutil"

	"github.com/user/numb/utils"
)

func GetConfig() NmbConfig {
	var nmbConfig NmbConfig
	raw, err := ioutil.ReadFile("nmb.json")
	utils.Check(err)
	utils.Check(json.Unmarshal(raw, &nmbConfig))
	return nmbConfig
}
