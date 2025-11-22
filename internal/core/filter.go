package core

import (
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// Filter handles file inclusion and exclusion logic
type Filter struct {
	extensions   map[string]struct{}
	matchers     []*ignore.GitIgnore
	excludeTests bool
}

// NewFilter creates a new Filter
func NewFilter(extensions []string, userPatterns []string, excludeTests bool) *Filter {
	extMap := make(map[string]struct{})
	for _, ext := range extensions {
		cleanExt := strings.TrimPrefix(ext, ".")
		extMap[cleanExt] = struct{}{}
	}

	var matchers []*ignore.GitIgnore

	// 1. System Constraints (Hard Blocks: Binaries, Git internals)
	// These are technically unreadable or dangerous to cat.
	systemIgnores := []string{
		".git",
		".DS_Store",
		"*.exe",
		"*.dll",
		"*.so",
		"*.dylib",
		"__pycache__",
		".venv",
		"venv",
		"node_modules",
		"target",
		"dist",
		"build",
		"*.log",
		"*.swp",
		".idea",
		".vscode",
		"vendor",
	}
	matchers = append(matchers, ignore.CompileIgnoreLines(systemIgnores...))

	// 2. Domain Opinions (Soft Blocks: Lockfiles, Noise)
	// These are text files we usually don't want, but might need.
	noiseIgnores := []string{
		"package-lock.json",
		"yarn.lock",
		"pnpm-lock.yaml",
		"go.sum",
		"Cargo.lock",
		"*.svg",
		"*.png",
		"*.jpg",
		"*.ico",
		"*.min.js",
		"*.min.css",
		"*.map",
	}

	// LOGIC: Only add a noise ignore if the user did NOT explicitly request it.
	var activeNoise []string
	for _, pattern := range noiseIgnores {
		// Heuristic: Check if the pattern's extension is in the requested list
		patternExt := strings.TrimPrefix(filepath.Ext(pattern), ".")
		// If it's a file like "go.sum", ext is "sum".
		// If it's "*.svg", ext is "svg".
		if _, requested := extMap[patternExt]; !requested {
			activeNoise = append(activeNoise, pattern)
		}
	}
	matchers = append(matchers, ignore.CompileIgnoreLines(activeNoise...))

	// 3. User Patterns & .gitignore (Standard behavior)
	// Add defaultIgnores that were previously there if not covered?
	// The previous list had: .git, node_modules, __pycache__, .venv, venv, target, dist, build, *.log, *.lock, *.swp, .DS_Store
	// dev tools: .idea, .vscode, vendor
	// cost saving: locks, media, mins.
	// I think I covered most.
	
	allPatterns := append([]string{}, userPatterns...)
	// We don't need to append defaultIgnores again as we split them.
	
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
		extensions:   extMap,
		matchers:     matchers,
		excludeTests: excludeTests,
	}
}

// HasValidExtension checks if the filename has a valid extension
func (f *Filter) HasValidExtension(filename string) bool {
	ext := filepath.Ext(filename)
	ext = strings.TrimPrefix(ext, ".")
	_, ok := f.extensions[ext]
	return ok
}

// IsTestFile checks if the file is a test file based on common conventions
func (f *Filter) IsTestFile(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	return strings.HasSuffix(base, "_test.go") || // Go
		strings.HasSuffix(base, ".test.js") || // JS/TS
		strings.HasSuffix(base, ".spec.js") || // JS/TS
		strings.HasSuffix(base, ".test.ts") || // JS/TS
		strings.HasSuffix(base, ".spec.ts") || // JS/TS
		strings.HasPrefix(base, "test_") // Python
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

	// NEW: Test Check
	if f.excludeTests && f.IsTestFile(path) {
		return false
	}

	// 2. Check extension inclusion
	if f.HasValidExtension(path) {
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
