package main

type LangProperties struct {
	ExecuteRule ExecuteRule `bson:"execute_rule" json:"execute_rule"`
	CompileRule CompileRule `bson:"compile_rule" json:"compile_rule"`
}

type ExecuteRule struct {
	CommandLine string `bson:"command_line" json:"command_line"`
	MaxMemory 	int64  `bson:"max_memory" json:"max_memory"`
	MaxCpuTime  int64  `bson:"max_cpu_time" json:"max_cpu_time"`
	MaxFileSize int64  `bson:"max_file_size" json:"max_file_size"`
}

type CompileRule struct {
	CommandLine string `bson:"command_line" json:"command_line"`
}
