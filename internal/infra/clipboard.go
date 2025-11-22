package infra

import (
	"github.com/atotto/clipboard"
)

// Clipboard handles clipboard operations
type Clipboard struct{}

// NewClipboard creates a new Clipboard instance
func NewClipboard() *Clipboard {
	return &Clipboard{}
}

// WriteAll writes the string to the system clipboard
func (c *Clipboard) WriteAll(text string) error {
	return clipboard.WriteAll(text)
}
