package static

import (
	"net/http"

	"github.com/renevo/gateway/env"
	"github.com/renevo/gateway/logging"
)

const (
	envBypassMemory  = "GATEWAY_SITE_MEMORY_FILE_DISABLE"
	envMaxMemorySize = "GATEWAY_SITE_MEMORY_FILE_MAX_SIZE"
)

// Site represents the gateway static hosting
type Site struct {
	handler http.Handler
	fs      http.FileSystem
}

// New creates a new static.Site
//
// By default this will load all files in the specified path less than 2mb into memory to serve without file IO
func New(path string) *Site {
	if env.Bool(envBypassMemory) {
		fs := http.Dir(path)
		return &Site{
			handler: http.FileServer(fs),
			fs:      fs,
		}
	}

	// create and wire up our memory file system
	fs, err := dir(path)

	if err != nil {
		logging.Errorf("Failed to read site path %q: %v", path, err)
	}

	return &Site{
		handler: http.FileServer(fs),
		fs:      fs,
	}
}

// ServeHTTP is the HTTP handler for the static web site
func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: default document (override the base code)
	// TODO: spa mode, if not found locally, serve the default document
	//		 this might require storing directories in the fs as well
	s.handler.ServeHTTP(w, r)
	return
}
