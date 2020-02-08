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

func compiler(req *CompileRequest) {
	// 언어에 따라 소스코드 파일 확장자 및 파일 이름 구성
	ext, err := languageExtension(req.LangType)
	if err != nil {
		log.Fatal(err)
	}
	filename := Md5HashGen("code") + ext

	// 폴더 생성
	path, err := MakePathDir(filepath.Join(BaseDirPath, Md5HashGen("test")))
	if err != nil {
		log.Fatal(err)
	}

	fullpath := filepath.Join(path, filename)
	fd, err := os.OpenFile(fullpath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0644))
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	w := bufio.NewWriter(fd)
	w.WriteString(req.SourceCode)
	w.Flush()
}

func languageExtension(langtype int) (string, error) {
	var ext string
	switch langtype {
	case C:
		ext = ".c"
		break
	case CXX:
		ext = ".cpp"
		break
	case JAVA:
		ext = ".java"
		break
	case PYTHON2:
	case PYTHON3:
		ext = ".py"
		break
	case GOLANG:
		ext = ".go"
		break
	case RUST:
		ext = ".rs"
		break

	default:
		return "", errors.New("Does not support language type")
	}

	return ext, nil
}
