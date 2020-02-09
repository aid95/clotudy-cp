package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
)

// 언어 종류에 대한 목록
const (
	C = iota
	CXX
	JAVA
	PYTHON2
	PYTHON3
	GOLANG
	RUST
)

// CompileRequest 컴파일 정보를 위한 구조체
type CompileRequest struct {
	ID             bson.ObjectId  `bson:"_id" json:"id"`
	CreatedAt      time.Time      `bson:"created_at" json:"created_at"`
	CableID        bson.ObjectId  `bson:"request_id" json:"request_id"`
	SourceCode     string         `bson:"src" json:"src"`
	SourceType     int            `bson:"type" json:"type"`
	LangProperties LangProperties `bson:"lang_properties" json:"lang_properties"`
}

// FieldMap 웹소켓으로 보내진 데이터를 CompileRequest 구조체의 요소와 맵핑
func (c *CompileRequest) FieldMap(r *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&c.LangProperties.CompileRule.SourceCode: "src",
		&c.LangProperties.CompileRule.LangType:   "type",
	}
}

func (c *CompileRequest) create() {
	c.ID = bson.NewObjectId()
	c.CreatedAt = time.Now()
	c.LangProperties.init()
}

func (c *CompileRequest) CompileAndRun() *ExecuteResponse {
	// 컴파일할 소스코드를 파일에 작성.
	fd, err := os.OpenFile(c.LangProperties.SourcePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer fd.Close()

	w := bufio.NewWriter(fd)
	if _, err := w.WriteString(c.SourceCode); err != nil {
		return nil
	}
	if err := w.Flush(); err != nil {
		return nil
	}

	er := &ExecuteResponse{}
	_, er.CompileOut, er.CompileErr = c.LangProperties.CompileRule.Run()
	_, er.ExecuteOut, er.ExecuteErr = c.LangProperties.ExecuteRule.Run()
	os.RemoveAll(c.LangProperties.BasePath)

	return er
}
