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
	BinaryPath  string      `bson:"bin_path" json:"bin_path"`
}

func (l *LangProperties) init() error {
	filename := Md5HashGen("code")
	path, err := MakePathDir(filepath.Join(BaseDirPath, Md5HashGen("test")))
	if err != nil {
		return err
	}
	MakePathDir(filepath.Join(path, "src"))
	MakePathDir(filepath.Join(path, "bin"))

	// 컴파일 시 필요한 파일 이름, 확장자, 경로 등을 설정
	l.BasePath = path
	l.SourcePath = fmt.Sprintf("%s/src/%s", l.BasePath, filename)
	l.BinaryPath = fmt.Sprintf("%s/bin/%s", l.BasePath, filename)
	switch l.CompileRule.LangType {
	case C:
		l.SourcePath += ".c"
		l.CompileRule.Compiler = "/usr/bin/gcc"
		l.CompileRule.CompileOption = []string{l.SourcePath, "-o", l.BinaryPath, "-O2", "-Wall", "-lm", "-static", "-std=c11"}
		l.ExecuteRule.Cmd = l.BinaryPath
		break
	case CXX:
		l.SourcePath += ".cpp"
		l.CompileRule.Compiler = "/usr/bin/g++"
		l.CompileRule.CompileOption = []string{l.SourcePath, "-o", l.BinaryPath, "-O2", "-Wall", "-lm", "-static", "-std=gnu++98"}
		l.ExecuteRule.Cmd = l.BinaryPath
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
	Cmd        string   `bson:"command_line" json:"command_line"`
	CmdOption  []string `bson:"command_option" json:"command_option"`
	MaxMem     int64    `bson:"max_memory" json:"max_memory"`
	MaxCPUTime int64    `bson:"max_cpu_time" json:"max_cpu_time"`
}

func (e *ExecuteRule) Run() (error, string, string) {
	return RunCommandLine(e.Cmd, e.CmdOption)
}

// CompileRule 컴파일 조건 및 명령행을 위한 정보
type CompileRule struct {
	Compiler      string   `bson:"lang_compiler" json:"lang_compiler"`
	CompileOption []string `bson:"compile_arg" json:"compile_arg"`
	MaxFileSize   int64    `bson:"max_file_size" json:"max_file_size"`
	SourceCode    string   `bson:"source_code" json:"source_code"`
	LangType      int      `bson:"source_type" json:"source_type"`
}

func (c *CompileRule) Run() (error, string, string) {
	return RunCommandLine(c.Compiler, c.CompileOption)
}
