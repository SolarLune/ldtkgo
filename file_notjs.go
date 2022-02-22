//go:build !js

package ldtkgo

import (
	"os"
	"path/filepath"
)

// OpenFile opens a file and returns a stream for its data.
//
// The path parts should be separated with slash '/' on any environments.
//
// Note that this doesn't work on mobiles.
func OpenFile(path string) (FileReader, error) {
	return os.Open(filepath.FromSlash(path))
}
