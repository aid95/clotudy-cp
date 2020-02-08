package main

import (
	"net/http"
	"time"

	"github.com/mholt/binding"

	"gopkg.in/mgo.v2/bson"
)

// CompileRequest 컴파일 정보를 위한 구조체
type CompileRequest struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	SourceCode string        `bson:"source_code" json:"source_code"`
	LangType   int           `bson:"source_type" json:"source_type"`
	CableID    bson.ObjectId `bson:"request_id" json:"request_id"`
}

// FieldMap 웹소켓으로 보내진 데이터를 CompileRequest 구조체의 요소와 맵핑
func (c *CompileRequest) FieldMap(r *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&c.SourceCode: "src",
		&c.LangType:   "type",
	}
}

func (c *CompileRequest) create() {
	c.ID = bson.NewObjectId()
	c.CreatedAt = time.Now()
}
