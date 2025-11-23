package config

type Config struct {
	Extensions     []string `yaml:"extensions"`
	IgnorePatterns []string `yaml:"ignore"`
	Output         string   `yaml:"output"`
	IncludeTree    bool     `yaml:"tree"`
	UseXML         bool     `yaml:"xml"`
	PrintToStdout  bool     `yaml:"stdout"`
	ExcludeTests   bool     `yaml:"no_tests"`
}
