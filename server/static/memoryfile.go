package static

import (
	"bytes"
	"net/http"
	"os"
	"time"
)

type memoryFile struct {
	contents []byte
	size     int64
	name     string
	fsPath   string
	urlPath  string
	mime     string
	modified time.Time
}

func (m *memoryFile) Name() string {
	return m.name
}

func (m *memoryFile) Size() int64 {
	return m.size
}

func (m *memoryFile) Mode() os.FileMode {
	return os.FileMode(0444) // readonly
}

func (m *memoryFile) ModTime() time.Time {
	return m.modified
}

func (m *memoryFile) IsDir() bool {
	return false
}

func (m *memoryFile) Sys() interface{} {
	return nil
}

func (m *memoryFile) Open() http.File {
	// should probably pool this at some point to cool down on GC
	return &httpFile{bytes.NewReader(m.contents), m}
}

type httpFile struct {
	*bytes.Reader
	info os.FileInfo
}

func (*httpFile) Readdir(count int) ([]os.FileInfo, error) {
	return []os.FileInfo{}, nil
}

func (f *httpFile) Stat() (os.FileInfo, error) {
	return f.info, nil
}

func (f *httpFile) Close() error {
	// gives us an opportunity to push this back to the pool
	return nil
}
