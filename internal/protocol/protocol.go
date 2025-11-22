package protocol

import "fmt"

const (
	// MarkerMD is the Markdown header format
	MarkerMD = "### File: %s ###"
	// MarkerXML is the XML header format (simplified for regex matching logic)
	MarkerXMLStart = `<file path="%s">`
	MarkerXMLEnd   = `</file>`
)

// FormatHeaderMD returns the formatted markdown header
func FormatHeaderMD(path string) string {
	return fmt.Sprintf(MarkerMD, path)
}

// FormatHeaderXML returns the formatted XML header
func FormatHeaderXML(path string) string {
	return fmt.Sprintf(MarkerXMLStart, path)
}
