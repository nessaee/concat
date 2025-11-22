package transform

import (
	"regexp"
	"strings"
)

// Options holds the transformation settings
type Options struct {
	Compact      bool
	StripHeaders bool
}

// Transformer handles text processing
type Transformer struct {
	options      Options
	multiNewline *regexp.Regexp
	headerBlock  *regexp.Regexp
	headerLine   *regexp.Regexp
	headerHash   *regexp.Regexp
	reMd         *regexp.Regexp
	reXml        *regexp.Regexp
}

// NewTransformer creates a new Transformer instance
func NewTransformer(opts Options) *Transformer {
	return &Transformer{
		options: opts,
		// Handle Windows line endings by matching \r?
		multiNewline: regexp.MustCompile(`(\r?\n){3,}`),
		headerBlock:  regexp.MustCompile(`(?s)^\s*/\*.*?(Copyright|License).*?\*/\s*`),
		headerLine:   regexp.MustCompile(`(?s)^(//.*(Copyright|License).*\n)+`),
		headerHash:   regexp.MustCompile(`(?s)^(#.*(Copyright|License).*\n)+`),
		// Protocol-aware Regexes
		// Protocol: ### File: %s ###
		reMd: regexp.MustCompile(`### File: (.*?) ###\s*\r?\n`),
		// Protocol: <file path="%s">
		reXml: regexp.MustCompile(`<file path="(.*?)">\s*\r?\n`),
	}
}

// Process applies all configured transformations to the content
func (t *Transformer) Process(content string) string {
	// Normalize line endings for consistent processing (optional, but recommended)
	content = strings.ReplaceAll(content, "\r\n", "\n")

	if t.options.StripHeaders {
		content = t.stripLicense(content)
	}

	if t.options.Compact {
		content = t.removeExcessWhitespace(content)
	}

	return content
}

func (t *Transformer) removeExcessWhitespace(content string) string {
	return t.multiNewline.ReplaceAllString(content, "\n\n")
}

func (t *Transformer) stripLicense(content string) string {
	// Try Markdown Header first
	if t.reMd.MatchString(content) {
		return t.splitAndClean(content, t.reMd)
	}

	// Try XML Header
	if t.reXml.MatchString(content) {
		return t.splitAndClean(content, t.reXml)
	}

	// Fallback (e.g. single file input without header)
	return t.stripLicenseSingle(content)
}

func (t *Transformer) splitAndClean(content string, re *regexp.Regexp) string {
	parts := re.Split(content, -1)
	matches := re.FindAllString(content, -1)

	var res strings.Builder
	if len(parts) > 0 {
		res.WriteString(parts[0]) // Preamble
	}

	for i, match := range matches {
		res.WriteString(match)
		if i+1 < len(parts) {
			cleaned := t.stripLicenseSingle(parts[i+1])
			res.WriteString(cleaned)
		}
	}
	return res.String()
}

func (t *Transformer) stripLicenseSingle(content string) string {
	c := t.headerBlock.ReplaceAllString(content, "")
	c = t.headerLine.ReplaceAllString(c, "")
	c = t.headerHash.ReplaceAllString(c, "")
	return strings.TrimSpace(c)
}
