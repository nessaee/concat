package transform

import (
	"bufio"
	"fmt"
	"io"
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

// ProcessStream applies transformations to the content stream line-by-line
// This enables O(1) memory usage relative to file size.
func (t *Transformer) ProcessStream(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	// Increase buffer size for long lines if needed
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	var licenseBuffer []string
	var pendingNewlines int
	inLicenseBlock := true
	const MaxLicenseLines = 50

	flushLicenseBuffer := func() {
		if len(licenseBuffer) == 0 {
			return
		}
		// Join buffer
		content := strings.Join(licenseBuffer, "\n")

		// Strip License (only if enabled)
		if t.options.StripHeaders {
			content = t.stripLicenseSingle(content)
		}
		// Compact (only if enabled)
		if t.options.Compact {
			content = t.removeExcessWhitespace(content)
		}

		// Write content if it's not empty
		if len(content) > 0 {
			fmt.Fprintln(w, content)
		}

		licenseBuffer = nil
		inLicenseBlock = false
		pendingNewlines = 0
	}

	for scanner.Scan() {
		line := scanner.Text()

		// Check for Header using fast string checks
		// MD: ### File: ... ###
		// XML: <file path="...">
		isHeader := false
		if strings.HasPrefix(line, "### File: ") && strings.HasSuffix(line, " ###") {
			isHeader = true
		} else if strings.HasPrefix(line, "<file path=\"") && strings.HasSuffix(line, "\">") {
			isHeader = true
		} else if line == "</file>" { // Handle XML closing tag
			isHeader = true
		}

		if isHeader {
			// New file started or ended. Flush previous buffer.
			if inLicenseBlock {
				flushLicenseBuffer()
			}

			// Write the header line directly
			fmt.Fprintln(w, line)

			// Reset for new file (or next segment)
			inLicenseBlock = true
			licenseBuffer = nil
			pendingNewlines = 0
			continue
		}

		// Body processing: Buffer initial lines for license detection
		if inLicenseBlock {
			licenseBuffer = append(licenseBuffer, line)
			if len(licenseBuffer) >= MaxLicenseLines {
				flushLicenseBuffer()
			}
			continue
		}

		// Streaming Body (Past license block)
		if t.options.Compact {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				pendingNewlines++
			} else {
				if pendingNewlines > 0 {
					fmt.Fprintln(w) // First newline
					if pendingNewlines > 1 {
						fmt.Fprintln(w) // Second newline
					}
					pendingNewlines = 0
				}
				fmt.Fprintln(w, line)
			}
		} else {
			// Just pass through
			fmt.Fprintln(w, line)
		}
	}

	// Final flush
	if inLicenseBlock {
		flushLicenseBuffer()
	}

	return scanner.Err()
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
