package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gocaine/go-dart/server/autogen"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// RerouteToIndex returns a new middleware to handle not found paths (e.g. frontend routes) and route them to index
// The blacklistPrefixes prevent some path to be routed
func RerouteToIndex(blacklistPrefixes ...string) gin.HandlerFunc {
	fs := autogen.FS(false)
	return func(c *gin.Context) {
		w, r := c.Writer, c.Request
		path := r.URL.Path
		for _, val := range blacklistPrefixes {
			if strings.HasPrefix(path, val) {
				// this has to be ignored
				return
			}
		}
		f, err := fs.Open("/index.html")
		if err != nil {
			msg, code := "500 Internal Server Error", http.StatusInternalServerError
			http.Error(w, msg, code)
			return
		}
		defer f.Close()

		d, err := f.Stat()
		if err != nil {
			msg, code := "500 Internal Server Error", http.StatusInternalServerError
			http.Error(w, msg, code)
			return
		}
		sizeFunc := func() (int64, error) {
			return d.Size(), nil
		}
		code := http.StatusOK
		ctype := mime.TypeByExtension(filepath.Ext(d.Name()))
		w.Header().Set("Content-Type", ctype)
		size, err := sizeFunc()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var sendContent io.Reader = f
		w.WriteHeader(code)

		if r.Method != "HEAD" {
			io.CopyN(w, sendContent, size)
		}
		c.Abort()
	}
}

// ServeStatics serves statics file
func ServeStatics() gin.HandlerFunc {
	fs := autogen.FS(false)
	fileserver := http.FileServer(fs)
	return func(c *gin.Context) {
		if exists(fs, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func exists(fs http.FileSystem, filepath string) bool {
	if !strings.HasPrefix(filepath, "/") {
		filepath = "/" + filepath
	}
	return autogen.Exists(fs, filepath)
}
