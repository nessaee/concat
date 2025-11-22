package core

import (
	"os"
	"path/filepath"
	"strings"
)

// TreeGenerator generates a directory tree structure
type TreeGenerator struct {
	filter *Filter
}

// NewTreeGenerator creates a new TreeGenerator
func NewTreeGenerator(filter *Filter) *TreeGenerator {
	return &TreeGenerator{filter: filter}
}

// Generate returns the tree structure as a string
func (t *TreeGenerator) Generate(root string) (string, error) {
	var sb strings.Builder
	sb.WriteString("### Directory Structure ###\n")
	sb.WriteString(".\n")

	// Use recursive generation
	treeStr, err := t.generateRecursive(root, "")
	if err != nil {
		return "", err
	}
	sb.WriteString(treeStr)
	return sb.String(), nil
}

func (t *TreeGenerator) generateRecursive(dir string, prefix string) (string, error) {
	var sb strings.Builder

	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// Filter entries first to know which is last
	var filtered []os.DirEntry
	for _, e := range entries {
		path := filepath.Join(dir, e.Name())

		// Handle relative path for filter if possible, or assume relative execution
		// If dir is ".", path is "foo".
		// If dir is "foo", path is "foo/bar".
		// This works for filter matchers usually.

		if t.filter.IsIgnored(path, e.IsDir()) {
			continue
		}
		
		// NEW: Tree Pruning - Only show files that match requested extensions
		if !e.IsDir() {
			if !t.filter.HasValidExtension(e.Name()) {
				continue
			}
		}

		filtered = append(filtered, e)
	}

	for i, e := range filtered {
		isLast := i == len(filtered)-1
		connector := "├── "
		newPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			newPrefix = prefix + "    "
		}

		sb.WriteString(prefix + connector + e.Name() + "\n")

		if e.IsDir() {
			path := filepath.Join(dir, e.Name())
			subTree, err := t.generateRecursive(path, newPrefix)
			if err != nil {
				return "", err
			}
			sb.WriteString(subTree)
		}
	}

	return sb.String(), nil
}
