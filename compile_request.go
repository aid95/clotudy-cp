package main

import (
	"github.com/mholt/binding"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// CompileRequest 컴파일 정보를 위한 구조체
type CompileRequest struct {
	// create() initialize
	ID         bson.ObjectId `bson:"_id" json:"id"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	// FieldMap mapping
	SourceCode string        `bson:"source_code" json:"source_code"`
	SourceType string        `bson:"source_type" json:"source_type"`
	// uninitialized
	CableID bson.ObjectId `bson:"request_id" json:"request_id"`
}

func (c *CompileRequest) FieldMap(r *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&c.SourceCode: "src",
		&c.SourceType: "type",
	}
}

func (c *CompileRequest) create() {
	c.ID = bson.NewObjectId()
	c.CreatedAt = time.Now()
}
