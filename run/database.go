// this file defines the database schema
package run

import (
	"time"

	"github.com/libgit2/git2go"

	"gopkg.in/mgo.v2/bson"
)

type Schema struct {
	ID                bson.ObjectId `bson:"_id,omitempty"`
	AbstractGraph     string        `json:"abstract" bson:"abstract"`
	ConcreteGraph     string        `json:"concrete" bson:"concrete"`
	Code              string        `json:"code" bson:"code"`
	StateDictFilename string        `json:"stateDictFilename" bson:"stateDictFilename"`
	Params            string        `json:"params" bson:"params"`
	Test              *TestInfo     `json:"test" bson:"test"`
	Timestamp         time.Time     `json:"timestamp" bson:"timestamp"`
	Versioning        *git.Oid      `json:"versioning" bson:"versioning"`
}

type TestInfo struct {
	Accuracy float64
	Extra    map[string]interface{}
	// TODO: more test records
}
