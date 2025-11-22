package core

import (
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// Filter handles file inclusion and exclusion logic
type Filter struct {
	extensions map[string]struct{}
	matchers   []*ignore.GitIgnore
}

// NewFilter creates a new Filter
func NewFilter(extensions []string, userPatterns []string) *Filter {
	extMap := make(map[string]struct{})
	for _, ext := range extensions {
		cleanExt := strings.TrimPrefix(ext, ".")
		extMap[cleanExt] = struct{}{}
	}

	var matchers []*ignore.GitIgnore

	// 1. Default + User Patterns
	defaultIgnores := []string{
		".git",
		"node_modules",
		"__pycache__",
		".venv",
		"venv",
		"target",
		"dist",
		"build",
		"*.log",
		"*.lock",
		"*.swp",
		".DS_Store",
	}
	allPatterns := append(defaultIgnores, userPatterns...)
	m1 := ignore.CompileIgnoreLines(allPatterns...)
	matchers = append(matchers, m1)

	// 2. .gitignore if exists
	if _, err := os.Stat(".gitignore"); err == nil {
		m2, err := ignore.CompileIgnoreFile(".gitignore")
		if err == nil {
			matchers = append(matchers, m2)
		}
	}

	return &Filter{
		extensions: extMap,
		matchers:   matchers,
	}
}

// ShouldProcess returns true if the file should be processed (included)
// path should be relative to the root
func (f *Filter) ShouldProcess(path string, isDir bool) bool {
	if f.IsIgnored(path, isDir) {
		return false
	}

	if isDir {
		return true
	}

	// 2. Check extension inclusion
	ext := filepath.Ext(path)
	ext = strings.TrimPrefix(ext, ".")

	if _, ok := f.extensions[ext]; ok {
		return true
	}

	return false
}

// IsIgnored returns true if the path matches any ignore pattern
func (f *Filter) IsIgnored(path string, isDir bool) bool {
	for _, m := range f.matchers {
		if m.MatchesPath(path) {
			return true
		}
		if isDir {
			if m.MatchesPath(path + "/") {
				return true
			}
		}
	}
	return false
}
