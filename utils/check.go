package utils

import "log"

var env = "development"

func Check(err error) {
	if err != nil {
		if env == "development" {
			panic(err)
		} else if env == "production" {
			log.Fatal(err)
		}
	}
}
