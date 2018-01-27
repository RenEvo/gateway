package static

import (
	"net/http"

	"github.com/renevo/gateway/logging"
)

// Site represents the gateway static hosting
type Site struct {
	files http.Handler
}

type siteFS struct {
	inner http.FileSystem
}

func (fs *siteFS) Open(name string) (http.File, error) {
	f, err := fs.inner.Open(name)

	logging.Infof("OpenFile: %q: %v", name, err)
	return f, err
}

func New(path string) *Site {
	fs := &siteFS{http.Dir(path)}

	return &Site{
		files: http.FileServer(fs),
	}
}

func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.files.ServeHTTP(w, r)
	return
}
