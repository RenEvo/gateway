package static

import (
	"bytes"
	"os"

	"github.com/renevo/gateway/logging"
)

type httpFile struct {
	*bytes.Reader
	info os.FileInfo
}

func (*httpFile) Readdir(count int) ([]os.FileInfo, error) {
	// TODO: Need to figure out what to do with this
	// I don't really have a use case for listing files in a directory for this application
	logging.Debugf("Read Directory; %d", count)
	return []os.FileInfo{}, nil
}

func (f *httpFile) Stat() (os.FileInfo, error) {
	return f.info, nil
}

func (f *httpFile) Close() error {
	// gives us an opportunity to push this back to the pool
	return nil
}
