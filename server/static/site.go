package static

import (
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
}

type siteFS struct {
	http.FileSystem
	root  string
	files sync.Map
}

func newFS(path string) *siteFS {
	absPath, _ := filepath.Abs(path)

	return &siteFS{
		FileSystem: http.Dir(path),
		root:       absPath,
	}
}

func (fs *siteFS) add(fsPath string, info os.FileInfo) {
	absPath, _ := filepath.Abs(fsPath)
	relPath, _ := filepath.Rel(fs.root, absPath)
	urlPath := "/" + strings.Replace(relPath, "\\", "/", -1)

	if info.IsDir() {
		logging.Debugf("Path: %q;Abs: %q; Rel: %q; Url: %q; Directory: %v", fsPath, absPath, relPath, urlPath, info.Name())
		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(info.Name()))
	logging.Debugf("Path: %q; Abs: %q; Rel: %q;Url: %q; Name: %s; Size: %d; Mime Type: %s", fsPath, absPath, relPath, urlPath, info.Name(), info.Size(), mimeType)

	// specifically don't large files (2mb)
	maxSize := env.Int64(envMaxMemorySize)
	if maxSize == 0 {
		maxSize = 2 * 1024 * 1024
	}

	if info.Size() > maxSize {
		return
	}

	// read it
	contents, err := ioutil.ReadFile(absPath)
	if err != nil {
		logging.Errorf("Failed to read file: %q; %v", absPath, err)
	}

	// store it for later
	fs.files.Store(urlPath, &memoryFile{
		contents: contents,
		fsPath:   absPath,
		mime:     mimeType,
		modified: info.ModTime(),
		name:     info.Name(),
		size:     info.Size(),
		urlPath:  urlPath,
	})
}

func (fs *siteFS) lookup(urlPath string) http.File {
	f, found := fs.files.Load(urlPath)
	if !found {
		return nil
	}

	return f.(*memoryFile).Open()
}

func (fs *siteFS) Open(name string) (http.File, error) {
	f := fs.lookup(name)
	if f != nil {
		logging.Debugf("OpenMemoryFile: %q", name)
		return f, nil
	}

	f, err := fs.FileSystem.Open(name)

	logging.Debugf("OpenFile: %q: %v", name, err)

	return f, err
}

// New creates a new static.Site
//
// By default this will load all files in the specified path less than 2mb into memory to serve without file IO
func New(path string) *Site {
	fs := newFS(path)

	if !env.Bool(envBypassMemory) {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				logging.Errorf("error with file %q: %v", path, err)
				return nil
			}
			fs.add(path, info)
			return nil
		})

		if err != nil {
			logging.Errorf("Failed to read site path %q: %v", path, err)
		}
	}

	return &Site{
		handler: http.FileServer(fs),
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
