package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
)

// CompileResult 컴파일 결과에 대한 정보
type CompileResult struct {
	src           string
	err           string
	compiletime   int32
	langtype      int32
	ExecuteResult ExecuteResult
}

// ExecuteResult 실행 결과에 대한 정보
type ExecuteResult struct {
	in      string
	out     string
	cputime int32
	memsize int32
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
