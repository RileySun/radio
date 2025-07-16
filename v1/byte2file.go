//AUTHOR/Saviour: github.com/Deng-Xian-Sheng/ (https://github.com/faiface/beep/issues/123)
package radio

import (
	"io"
	"io/fs"
	"reflect"
	"time"
)

// A file is a single file in the FS.
// It implements fs.FileInfo and fs.DirEntry.
type file struct {
	name string
	data string
	hash [16]byte // truncated SHA256 hash
}

func (f *file) Name() string               { return f.name }
func (f *file) Size() int64                { return int64(len(f.data)) }
func (f *file) ModTime() time.Time         { return time.Time{} }
func (f *file) IsDir() bool                { return false }
func (f *file) Sys() any                   { return nil }
func (f *file) Type() fs.FileMode          { return f.Mode().Type() }
func (f *file) Info() (fs.FileInfo, error) { return f, nil }

func (f *file) Mode() fs.FileMode {
	if f.IsDir() {
		return fs.ModeDir | 0555
	}
	return 0444
}

// An openFile is a regular file open for reading.
type openFile struct {
	f      *file // the file itself
	offset int64 // current read offset
}

func (f *openFile) Close() error               { return nil }
func (f *openFile) Stat() (fs.FileInfo, error) { return f.f, nil }

// Read reads up to len(b) bytes from the File and stores them in b.
// It returns the number of bytes read and any error encountered.
// At end of file, Read returns 0, io.EOF.
func (f *openFile) Read(b []byte) (int, error) {
	if f.offset >= int64(len(f.f.data)) {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: f.f.name, Err: fs.ErrInvalid}
	}
	n := copy(b, f.f.data[f.offset:])
	f.offset += int64(n)
	return n, nil
}

var (
	_ io.Seeker   = (*openFile)(nil)
	_ fs.FileInfo = (*file)(nil)
)

func (f *openFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		// offset += 0
	case 1:
		offset += f.offset
	case 2:
		offset += int64(len(f.f.data))
	}
	if offset < 0 || offset > int64(len(f.f.data)) {
		return 0, &fs.PathError{Op: "seek", Path: f.f.name, Err: fs.ErrInvalid}
	}
	f.offset = offset
	return offset, nil
}

func (f *openFile) Write(b []byte) (int, error) {
	if reflect.ValueOf(f).IsNil() {
		return 0, &fs.PathError{Op: "write", Path: f.f.name, Err: fs.ErrInvalid}
	}
	f.f.data += string(b)
	return len(b), nil
}

func NewVirtualFile() *openFile {
	return &openFile{
		f:      &file{},
		offset: 0,
	}
}