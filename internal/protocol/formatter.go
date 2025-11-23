package protocol

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Formatter defines the interface for output formatting
type Formatter interface {
	WriteHeader(w io.Writer, path string)
	WriteFooter(w io.Writer)
}

// MarkdownFormatter implements Formatter for Markdown output
type MarkdownFormatter struct{}

func (f *MarkdownFormatter) WriteHeader(w io.Writer, path string) {
	fmt.Fprintf(w, MarkerMD+"\n", path)
}

func (f *MarkdownFormatter) WriteFooter(w io.Writer) {
	fmt.Fprint(w, "\n\n---\n\n")
}

// XMLFormatter implements Formatter for XML output
type XMLFormatter struct{}

func (f *XMLFormatter) WriteHeader(w io.Writer, path string) {
	// Manually construct the tag to ensure the path attribute is escaped
	// We can't use xml.Encoder easily here because we aren't encoding the whole element structure at once,
	// just the opening tag.
	// <file path="...">
	fmt.Fprint(w, "<file path=\"")
	if err := xml.EscapeText(w, []byte(path)); err != nil {
		// Fallback to raw string if escaping fails (should not happen for strings)
		fmt.Fprint(w, path)
	}
	fmt.Fprint(w, "\">\n")
}

func (f *XMLFormatter) WriteFooter(w io.Writer) {
	fmt.Fprint(w, "\n</file>\n")
}
