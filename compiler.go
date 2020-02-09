package main

import (
	"bufio"
	"github.com/mholt/binding"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// CompileRequest 컴파일 정보를 위한 구조체
type CompileRequest struct {
	ID             bson.ObjectId  `bson:"_id" json:"id"`
	CreatedAt      time.Time      `bson:"created_at" json:"created_at"`
	CableID        bson.ObjectId  `bson:"request_id" json:"request_id"`
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

func compiler(cr *CompileRequest, s *Service) error {
	// 컴파일할 소스코드를 파일에 작성.
	p := filepath.Join(cr.LangProperties.BasePath, cr.LangProperties.SourcePath)
	fd, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer fd.Close()

	w := bufio.NewWriter(fd)
	if _, err := w.WriteString(cr.LangProperties.CompileRule.SourceCode); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}

	res := &ExecuteResponse{}
	if out, err := exec.Command(cr.LangProperties.CompileRule.Cmd).Output(); err != nil {
		log.Fatal(err)
		return err
	} else {
		res.CompileOut = string(out)
	}

	if out, err := exec.Command(cr.LangProperties.ExecuteRule.Cmd).Output(); err != nil {
		log.Fatal(err)
		return err
	} else {
		res.ExecuteOut = string(out)
	}
	s.Send <- res

	return nil
}
