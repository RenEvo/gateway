package static

import "testing"

func TestIsDir(t *testing.T) {
	f := &fileInfo{
		directory: true,
	}

	if !f.IsDir() {
		t.Errorf("fileInfo did not return IsDir()")
	}

	if !f.Mode().IsDir() {
		t.Errorf("fileInfo.Mode() did not return IsDir()")
	}
}
