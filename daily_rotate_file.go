package log

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

// File describes a file that gets rotated daily
type File struct {
	sync.Mutex
	pathFormat string

	// info about currently opened file
	day      int
	path     string
	file     *os.File
	onRotate func(path string)
}

func (f *File) close() error {
	var err error
	if f.file != nil {
		err = f.file.Close()
		f.file = nil
	}
	return err
}

func (f *File) open() error {
	t := time.Now().UTC()
	f.path = t.Format(f.pathFormat)
	f.day = t.YearDay()

	// we can't assume that the dir for the file already exists
	dir := filepath.Dir(f.path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	flag := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	f.file, err = os.OpenFile(f.path, flag, 0644)
	return err
}

// rotate on new day
func (f *File) reopenIfNeeded() error {
	t := time.Now().UTC()
	if t.YearDay() == f.day {
		return nil
	}
	err := f.close()
	if err != nil {
		return err
	}
	if f.onRotate != nil {
		f.onRotate(f.path)
	}
	return f.open()
}

// NewFile opens a new log file (creates if doesn't exist, will append if exists)
func NewFile(pathFormat string, onRotate func(path string)) (*File, error) {
	res := &File{
		pathFormat: pathFormat,
		onRotate:   onRotate,
	}
	if err := res.open(); err != nil {
		return nil, err
	}
	return res, nil
}

// Close closes the file
func (f *File) Close() error {
	f.Lock()
	defer f.Unlock()
	return f.close()
}

// Write writes data to a file
func (f *File) Write(d []byte) (int, error) {
	f.Lock()
	defer f.Unlock()
	err := f.reopenIfNeeded()
	if err != nil {
		return 0, err
	}
	return f.file.Write(d)
}

// Flush flushes the file
func (f *File) Flush() error {
	f.Lock()
	defer f.Unlock()
	return f.file.Sync()
}
