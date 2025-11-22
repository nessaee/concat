package core

import (
	"bytes"
)

// isBinary checks if the data is likely binary by looking for null bytes
// in the first 8000 bytes (similar to git's heuristic).
func isBinary(content []byte) bool {
	limit := 8000
	if len(content) < limit {
		limit = len(content)
	}
	return bytes.IndexByte(content[:limit], 0) != -1
}
