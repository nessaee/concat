package config

type Config struct {
	Extensions     []string
	IgnorePatterns []string
	Output         string
	IncludeTree    bool
	UseXML         bool
	PrintToStdout  bool
	ExcludeTests   bool
}
