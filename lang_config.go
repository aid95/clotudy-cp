package main

// LangProperties 언어별 실행 및 컴팡ㄹ 정보를 위한 정보
type LangProperties struct {
	ExecuteRule ExecuteRule `bson:"execute_rule" json:"execute_rule"`
	CompileRule CompileRule `bson:"compile_rule" json:"compile_rule"`
	BasePath    string      `bson:"base_path" json:"base_path"`
	SourcePath  string      `bson:"src_path" json:"src_path"`
	BinaryPath  string      `bson:"bin_path" json:"bin_path"`
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
