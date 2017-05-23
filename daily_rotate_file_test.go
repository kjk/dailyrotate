package dailyrotate

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	err := os.RemoveAll("test_dir")
	assert.Nil(t, err)
	pathFormat := filepath.Join("test_dir", "second", "2006-01-02.txt")
	pathExp := time.Now().UTC().Format(pathFormat)
	f, err := NewFile(pathFormat, nil)
	assert.Nil(t, err)
	assert.NotNil(t, f)
	n, err := io.WriteString(f, "hello\n")
	assert.Nil(t, err)
	assert.Equal(t, 0, int(f.lastWriteCurrPos))
	assert.Equal(t, n, 6)
	assert.Equal(t, pathExp, f.path)
	n, err = f.Write([]byte("bar\n"))
	assert.Nil(t, err)
	assert.Equal(t, 6, int(f.lastWriteCurrPos))
	assert.Equal(t, n, 4)
	err = f.Close()
	assert.Nil(t, err)

	d, err := ioutil.ReadFile(pathExp)
	assert.Nil(t, err)
	assert.Equal(t, string(d), "hello\nbar\n")

	path, off, n, err := f.Write2([]byte("and more\n"), true)
	assert.Nil(t, err)

	assert.Equal(t, len(d), int(off))

	assert.Equal(t, 9, n)
	assert.Equal(t, pathExp, path)

	err = f.Close()
	assert.Nil(t, err)
	d, err = ioutil.ReadFile(pathExp)
	assert.Nil(t, err)
	assert.Equal(t, string(d), "hello\nbar\nand more\n")

	err = os.RemoveAll("test_dir")
	assert.Nil(t, err)
}
