package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
		&c.SourceCode: "src",
		&c.SourceType: "type",
	}
}

func (c *CompileRequest) create() {
	c.ID = bson.NewObjectId()
	c.CreatedAt = time.Now()
	c.init()
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
	// os.RemoveAll(c.LangProperties.BasePath)

	return er
}

func (c *CompileRequest) init() error {
	filename := Md5HashGen("code")
	path, err := MakePathDir(filepath.Join(BaseDirPath, Md5HashGen("test")))
	if err != nil {
		return err
	}
	MakePathDir(filepath.Join(path, "src"))
	MakePathDir(filepath.Join(path, "bin"))

	// 컴파일 시 필요한 파일 이름, 확장자, 경로 등을 설정
	// 실행 시간, 메모리 제한을 위한 timeout 설치
	// COMMAND:
	//  curl https://raw.githubusercontent.com/pshved/timeout/master/timeout | \
	//  sudo tee /usr/bin/timeout && sudo chmod 755 /usr/bin/timeout"
	c.LangProperties.BasePath = path
	c.LangProperties.SourcePath = fmt.Sprintf("%s/src/%s", c.LangProperties.BasePath, filename)
	c.LangProperties.BinaryPath = fmt.Sprintf("%s/bin/%s", c.LangProperties.BasePath, filename)
	switch c.SourceType {
	case C:
		c.LangProperties.SourcePath += ".c"

		c.LangProperties.CompileRule.Compiler = "/usr/bin/gcc"
		c.LangProperties.CompileRule.CompileOption = []string{c.LangProperties.SourcePath, "-o", c.LangProperties.BinaryPath, "-O2", "-Wall", "-lm", "-static", "-std=c11"}
		c.LangProperties.ExecuteRule.Cmd = "/usr/bin/timeout"
		c.LangProperties.ExecuteRule.CmdOption = []string{"-m", "500", "-t", "3", c.LangProperties.BinaryPath}
		break
	case CXX:
		c.LangProperties.SourcePath += ".cpp"

		c.LangProperties.CompileRule.Compiler = "/usr/bin/g++"
		c.LangProperties.CompileRule.CompileOption = []string{c.LangProperties.SourcePath, "-o", c.LangProperties.BinaryPath, "-O2", "-Wall", "-lm", "-static", "-std=gnu++98"}
		c.LangProperties.ExecuteRule.Cmd = "/usr/bin/timeout"
		c.LangProperties.ExecuteRule.CmdOption = []string{"-m", "500", "-t", "3", c.LangProperties.BinaryPath}
		break
	case JAVA:
		break
	case PYTHON2:
		c.LangProperties.SourcePath += ".py"

		c.LangProperties.CompileRule.Compiler = "/usr/bin/python"
		c.LangProperties.CompileRule.CompileOption = []string{"-c", fmt.Sprintf("\"import py_compile; py_compile.compile(r'%s')\"", c.LangProperties.SourcePath)}
		c.LangProperties.ExecuteRule.Cmd = "/usr/bin/python"
		c.LangProperties.ExecuteRule.CmdOption = []string{c.LangProperties.SourcePath}
		break
	case PYTHON3:
		c.LangProperties.SourcePath += ".py"

		c.LangProperties.CompileRule.Compiler = "/usr/bin/python3"
		c.LangProperties.CompileRule.CompileOption = []string{"-c", fmt.Sprintf("\"import py_compile; py_compile.compile(r'%s')\"", c.LangProperties.SourcePath)}
		c.LangProperties.ExecuteRule.Cmd = "/usr/bin/python3"
		c.LangProperties.ExecuteRule.CmdOption = []string{c.LangProperties.SourcePath}
		break
	case GOLANG:
		c.LangProperties.SourcePath += ".go"

		c.LangProperties.CompileRule.Compiler = "/usr/bin/go"
		c.LangProperties.CompileRule.CompileOption = []string{"-c", fmt.Sprintf("\"import py_compile; py_compile.compile(r'%s')\"", c.LangProperties.SourcePath)}
		c.LangProperties.ExecuteRule.Cmd = "/usr/bin/python3"
		c.LangProperties.ExecuteRule.CmdOption = []string{c.LangProperties.SourcePath}
		break
	case RUST:
		break
	default:
		return errors.New("Does not support language type.")
	}
	return nil
}
