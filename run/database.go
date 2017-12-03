// this file defines the database schema
package run

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Schema struct {
	ID                bson.ObjectId `bson:"_id,omitempty"`
	AbstractGraph     string        `json:"abstract" bson:"abstract"`
	ConcreteGraph     string        `json:"concrete" bson:"concrete"`
	Code              string        `json:"code" bson:"code"`
	StateDictFilename string        `json:"stateDictFilename" bson:"stateDictFilename"`
	Params            string        `json:"params" bson:"params"`
	Test              string     `json:"test" bson:"test"`
	Timestamp         time.Time     `json:"timestamp" bson:"timestamp"`
	Versioning        string      `json:"versioning" bson:"versioning"`
}

