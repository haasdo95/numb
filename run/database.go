// this file defines the database schema
package run

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// type Schema struct {
// 	ID                bson.ObjectId `bson:"_id,omitempty"`
// 	AbstractGraph     string
// 	ConcreteGraph     string
// 	StateDictFilename string
// 	Params            string
// 	Test              *TestInfo
// 	Timestamp         time.Time
// }

type Schema struct {
	ID                bson.ObjectId `bson:"_id,omitempty"`
	AbstractGraph     string        `json:"abstract" bson:"abstract"`
	ConcreteGraph     string        `json:"concrete" bson:"concrete"`
	StateDictFilename string        `json:"stateDictFilename" bson:"stateDictFilename"`
	Params            string        `json:"params" bson:"params"`
	Test              *TestInfo     `json:"test" bson:"test"`
	Timestamp         time.Time     `json:"timestamp" bson:"timestamp"`
}

type TestInfo struct {
	Accuracy float64
	Extra    map[string]interface{}
	// TODO: more test records
}
