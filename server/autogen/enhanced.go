package autogen

import (
	"net/http"
	"path"
)

// Exists checks if a given filePath is present in the given FS
func Exists(fs http.FileSystem, filePath string) bool {
	_, present := _escData[path.Clean(filePath)]
	return present
}
