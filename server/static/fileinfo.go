package static

import (
	"bytes"
	"net/http"
	"os"
	"time"

	"github.com/renevo/gateway/logging"
)

const (
	defaultFilePermissions = 0444
	directoryFileMode      = os.FileMode(os.ModeDir | defaultFilePermissions)
	defaultFileMode        = os.FileMode(defaultFilePermissions)
)

type fileInfo struct {
	directory bool
	contents  []byte
	size      int64
	name      string
	fsPath    string
	urlPath   string
	mime      string
	modified  time.Time
}

func (m *fileInfo) Name() string {
	return m.name
}

func (m *fileInfo) Size() int64 {
	return m.size
}

func (m *fileInfo) Mode() os.FileMode {
	if m.directory {
		return directoryFileMode
	}

	return defaultFileMode
}

func (m *fileInfo) ModTime() time.Time {
	return m.modified
}

func (m *fileInfo) IsDir() bool {
	return m.directory
}

func (m *fileInfo) Sys() interface{} {
	return nil
}

func (m *fileInfo) Open() (http.File, error) {
	logging.Debugf("Memory: Opening File: %q", m.name)

	if m.size > 0 && len(m.contents) == 0 {
		logging.Debugf("Memory: Opening File From FS: %s", m.fsPath)
		return os.Open(m.fsPath)
	}

	// should probably pool this at some point to cool down on GC
	return &httpFile{bytes.NewReader(m.contents), m}, nil
}
