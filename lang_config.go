package main

// LangProperties 언어별 실행 및 컴팡ㄹ 정보를 위한 정보
type LangProperties struct {
	ExecuteRule ExecuteRule `bson:"execute_rule" json:"execute_rule"`
	CompileRule CompileRule `bson:"compile_rule" json:"compile_rule"`
}

// ExecuteRule 실행을 위한 정보
type ExecuteRule struct {
	CommandLine string `bson:"command_line" json:"command_line"`
	MaxMemory   int64  `bson:"max_memory" json:"max_memory"`
	MaxCPUTime  int64  `bson:"max_cpu_time" json:"max_cpu_time"`
	MaxFileSize int64  `bson:"max_file_size" json:"max_file_size"`
}

// CompileRule 컴파일 조건 및 명령행을 위한 정보
type CompileRule struct {
	CommandLine string `bson:"command_line" json:"command_line"`
}
