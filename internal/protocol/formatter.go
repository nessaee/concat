package protocol

import (
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
	fmt.Fprintf(w, MarkerXMLStart+"\n", path)
}

func (f *XMLFormatter) WriteFooter(w io.Writer) {
	fmt.Fprintf(w, "\n"+MarkerXMLEnd+"\n")
}
