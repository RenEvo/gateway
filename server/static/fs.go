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

// TODO: Want to add a file watcher for the static to update if env is set (dev mode)
// How to detect file changes in Golang
// https://medium.com/@skdomino/watch-this-file-watching-in-go-5b5a247cf71f

type siteFS struct {
	http.FileSystem
	root  string
	files sync.Map
}

func dir(path string) (*siteFS, error) {
	absPath, _ := filepath.Abs(path)

	fs := &siteFS{
		FileSystem: http.Dir(path),
		root:       absPath,
	}

	// specifically don't large files (2mb default)
	maxSize := env.Int64(envMaxMemorySize)
	if maxSize == 0 {
		maxSize = 2 * 1024 * 1024
	}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logging.Errorf("error with file %q: %v", path, err)
			return err
		}

		return fs.add(maxSize, path, info)
	})

	if err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *siteFS) add(maxSize int64, fsPath string, info os.FileInfo) error {
	absPath, _ := filepath.Abs(fsPath)
	relPath, _ := filepath.Rel(fs.root, absPath)
	urlPath := "/" + strings.Replace(relPath, "\\", "/", -1)

	if info.IsDir() {
		logging.Debugf("Path: %q; Abs: %q; Rel: %q; Url: %q; Directory: %v", fsPath, absPath, relPath, urlPath, info.Name())
		fs.files.Store(urlPath, &fileInfo{
			directory: true,
			fsPath:    absPath,
			modified:  info.ModTime(),
			name:      info.Name(),
			urlPath:   urlPath,
		})
		return nil
	}

	mimeType := mime.TypeByExtension(filepath.Ext(info.Name()))
	logging.Debugf("Path: %q; Abs: %q; Rel: %q;Url: %q; Name: %s; Size: %d; Mime Type: %s", fsPath, absPath, relPath, urlPath, info.Name(), info.Size(), mimeType)

	// if the file is too big then only store the reference to the file
	if info.Size() > maxSize {
		fs.files.Store(urlPath, &fileInfo{
			fsPath:   absPath,
			mime:     mimeType,
			modified: info.ModTime(),
			name:     info.Name(),
			size:     info.Size(),
			urlPath:  urlPath,
		})
		return nil
	}

	// read it
	contents, err := ioutil.ReadFile(absPath)
	if err != nil {
		logging.Errorf("Failed to read file: %q; %v", absPath, err)
		return err
	}

	// store it for later
	fs.files.Store(urlPath, &fileInfo{
		contents: contents,
		fsPath:   absPath,
		mime:     mimeType,
		modified: info.ModTime(),
		name:     info.Name(),
		size:     info.Size(),
		urlPath:  urlPath,
	})

	return nil
}

func (fs *siteFS) lookup(urlPath string) *fileInfo {
	f, found := fs.files.Load(urlPath)
	if !found {
		return nil
	}

	return f.(*fileInfo)
}

func (fs *siteFS) Open(name string) (http.File, error) {
	mf := fs.lookup(name)
	if mf != nil {
		logging.Debugf("OpenMemoryFile: %q", name)
		return mf.Open()
	}

	f, err := fs.FileSystem.Open(name)

	logging.Debugf("OpenFile: %q: %v", name, err)

	return f, err
}
