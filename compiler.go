package main

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
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
}

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

func compiler(cr *CompileRequest, s *Service) (error, error) {
	// 컴파일할 소스코드를 파일에 작성.
	fd, err := os.OpenFile(cr.LangProperties.SourcePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer fd.Close()

	w := bufio.NewWriter(fd)
	if _, err := w.WriteString(cr.SourceCode); err != nil {
		return nil, err
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}

	var cstdout bytes.Buffer
	var cstderr bytes.Buffer
	res := &ExecuteResponse{}

	// Compile
	ccmd := exec.Command(cr.LangProperties.CompileRule.Compiler, cr.LangProperties.CompileRule.CompileOption...)
	ccmd.Stderr = &cstderr
	ccmd.Stdout = &cstdout
	compileErr := ccmd.Run()
	res.CompileErr = cstderr.String()
	res.CompileOut = cstdout.String()

	var estdout bytes.Buffer
	var estderr bytes.Buffer
	// Execute
	ecmd := exec.Command(cr.LangProperties.ExecuteRule.Cmd, cr.LangProperties.ExecuteRule.CmdOption)
	ecmd.Stderr = &estderr
	ecmd.Stdout = &estdout
	executeErr := ecmd.Run()
	res.ExecuteErr = estderr.String()
	res.ExecuteOut = estdout.String()

	s.Send <- res

	return compileErr, executeErr
}
