// this file defines the database schema
package run

import (
	"gopkg.in/mgo.v2/bson"
)

type Schema struct {
	ID                bson.ObjectId `bson:"_id,omitempty"`
	AbstractGraph     string
	ConcreteGraph     string
	StateDictFilename string
	Params            map[string]interface{}
	Test              TestInfo
}

type TestInfo struct {
	Accuracy float64
	Extra    map[string]interface{}
	// TODO: more test records
}
