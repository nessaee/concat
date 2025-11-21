package config

// Config holds the application configuration
type Config struct {
	Extensions     []string
	IgnorePatterns []string
	Output         string
	IncludeTree    bool
	ShowVersion    bool
}
