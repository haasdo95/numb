package bootstrap

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createNmbJSON() error {
	nmbjson := NmbConfig{
		Train: make([]string, 0),
		Test:  make([]string, 0),
	}

	jsonbytes, err := json.MarshalIndent(nmbjson, "", "\t")
	check(err)

	if _, err := os.Stat("nmb.json"); os.IsNotExist(err) { // nmb.json does not exist
		jsonFile, err := os.OpenFile("nmb.json", os.O_RDWR|os.O_CREATE, 0666)
		check(err)
		defer jsonFile.Close()

		_, err = jsonFile.Write(jsonbytes)
		check(err)

		return nil

	} else { // nmb.json already exists
		return errors.New("nmb.json already exists!")
	}

}

func createNmbDir() error {
	if _, err := os.Stat(".nmb"); os.IsNotExist(err) {
		check(os.Mkdir(".nmb", 0755))
		return nil
	} else {
		return errors.New(".nmb already exists!")
	}
}

func removeNmbJSON() error {
	err := os.Remove("nmb.json")
	if err != nil && os.IsNotExist(err) {
		return errors.New("nmb.json doesn't even exist")
	}
	return nil
}

func removeNmbDir(recursive bool) error {
	err := os.Remove(".nmb")
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New(".nmb doesn't even exist")
		} else {
			if recursive {
				check(os.RemoveAll(".nmb"))
			} else {
				return errors.New("set recursive as True to force removal")
			}
		}
	}
	return nil
}

// Deinit serves to remove both .nmb and nmb.json
func Deinit() (error, error) {
	errdir := removeNmbDir(true)
	if errdir != nil {
		fmt.Println(errdir.Error())
	}
	errjson := removeNmbJSON()
	if errjson != nil {
		fmt.Println(errjson.Error())
	}
	return errdir, errjson
}

// Init serves to remove both .nmb and nmb.json
func Init() (error, error) {
	errdir := createNmbDir()
	if errdir != nil {
		fmt.Println(errdir.Error())
	}
	errjson := createNmbJSON()
	if errjson != nil {
		fmt.Println(errjson.Error())
	}
	return errdir, errjson
}
