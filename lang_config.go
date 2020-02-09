package main

import (
	"errors"
	"fmt"
	"path/filepath"
)

// LangProperties 언어별 실행 및 컴팡ㄹ 정보를 위한 정보
type LangProperties struct {
	ExecuteRule ExecuteRule `bson:"execute_rule" json:"execute_rule"`
	CompileRule CompileRule `bson:"compile_rule" json:"compile_rule"`
	BasePath    string      `bson:"base_path" json:"base_path"`
	SourcePath  string      `bson:"src_path" json:"src_path"`
	BinaryPath	string		`bson:"bin_path" json:"bin_path"`
}

func (l *LangProperties) init() error {
	filename := Md5HashGen("code")
	path, err := MakePathDir(filepath.Join(BaseDirPath, Md5HashGen("test")))
	if err != nil {
		return err
	}

	// 컴파일 시 필요한 파일 이름, 확장자, 경로 등을 설정
	l.BasePath = path
	l.SourcePath = fmt.Sprintf("%s/src/%s", l.BasePath, filename)
	l.BinaryPath = fmt.Sprintf("%s/bin/%s", l.BasePath, filename)
	switch l.CompileRule.LangType {
	case C:
		l.SourcePath += ".c"
		l.CompileRule.Cmd = fmt.Sprintf("gcc %s -o %s -O2 -Wall -lm -static -std=c11", l.SourcePath, l.BinaryPath)
		break
	case CXX:
		l.SourcePath += ".cpp"
		l.CompileRule.Cmd = fmt.Sprintf("g++ %s -o %s -O2 -Wall -lm -static -std=gnu++98", l.SourcePath, l.BinaryPath)
		break
	case JAVA:
		break
	case PYTHON2:
	case PYTHON3:
		break
	case GOLANG:
		break
	case RUST:
		break
	default:
		return errors.New("Does not support language type.")
	}
	return nil
}

// ExecuteRule 실행을 위한 정보
type ExecuteRule struct {
	Cmd        string `bson:"command_line" json:"command_line"`
	MaxMem     int64  `bson:"max_memory" json:"max_memory"`
	MaxCPUTime int64  `bson:"max_cpu_time" json:"max_cpu_time"`
}

// CompileRule 컴파일 조건 및 명령행을 위한 정보
type CompileRule struct {
	Cmd         string `bson:"command_line" json:"command_line"`
	MaxFileSize int64  `bson:"max_file_size" json:"max_file_size"`
	SourceCode  string `bson:"source_code" json:"source_code"`
	LangType    int    `bson:"source_type" json:"source_type"`
}
