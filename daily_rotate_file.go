package dailyrotate

import (
	"io"
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
	onRotate func(path string, didRotate bool)

	// for tests only
	lastWriteCurrPos int64
}

func (f *File) close(didRotate bool) error {
	if f.file == nil {
		return nil
	}
	err := f.file.Close()
	f.file = nil
	if f.onRotate != nil && err != nil {
		f.onRotate(f.path, didRotate)
	}
	f.day = 0
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

	// would be easier to open with os.O_APPEND but Seek() doesn't work in that case
	flag := os.O_CREATE | os.O_WRONLY
	f.file, err = os.OpenFile(f.path, flag, 0644)
	if err != nil {
		return err
	}
	_, err = f.file.Seek(0, io.SeekEnd)
	return err
}

// rotate on new day
func (f *File) reopenIfNeeded() error {
	t := time.Now().UTC()
	if t.YearDay() == f.day {
		return nil
	}
	err := f.close(true)
	if err != nil {
		return err
	}
	return f.open()
}

// NewFile opens a new log file (creates if doesn't exist, will append if exists)
func NewFile(pathFormat string, onRotate func(path string, didRotate bool)) *File {
	return &File{
		pathFormat: pathFormat,
		onRotate:   onRotate,
	}
}

// Close closes the file
func (f *File) Close() error {
	f.Lock()
	defer f.Unlock()
	return f.close(false)
}

func (f *File) write(d []byte, flush bool) (int64, int, error) {
	err := f.reopenIfNeeded()
	if err != nil {
		return 0, 0, err
	}
	f.lastWriteCurrPos, err = f.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, 0, err
	}
	n, err := f.file.Write(d)
	if err != nil {
		return 0, n, err
	}
	if flush {
		err = f.file.Sync()
	}
	return f.lastWriteCurrPos, n, err
}

// Write writes data to a file
func (f *File) Write(d []byte) (int, error) {
	f.Lock()
	defer f.Unlock()
	_, n, err := f.write(d, false)
	return n, err
}

// Write2 writes data to a file, optionally flushes. To enable users to later
// seek to where the data was written, it returns name of the file where data
// was written, offset at which the data was written, number of bytes and error
func (f *File) Write2(d []byte, flush bool) (string, int64, int, error) {
	f.Lock()
	defer f.Unlock()
	currPos, n, err := f.write(d, flush)
	return f.path, currPos, n, err
}

// Flush flushes the file
func (f *File) Flush() error {
	f.Lock()
	defer f.Unlock()
	return f.file.Sync()
}
