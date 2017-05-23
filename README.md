# dailyrotate

Go library for a file whose that rotates daily.

An io.Writer implementation that writes to a file
and rotates daily (according to UTC() time).

## Usage

```go
pathFormat := filepath.Join("dir", "2006-01-02.txt")
func onClose(path string, didRotate bool) {
  fmt.Printf("we just closed a file '%s', didRotate: %v\n", path, didRotate)
}
w, err := dailyrotate.NewFile(pathFormat, onClose)
panicIfErr(err)
_, err = io.WriteString(w, "hello\n")
panicIfErr(err)
err = w.Close()
panicIfErr(err)
```

Given that files are rotated daily, you need to provide
a file name format which will be passed to time.Now().UTC().Format()
function. It should produce a unique file name each day.

You can also provide `onRotate` function which will be called
when the file is close (either because of rotation or regular Close()).

You can use that to e.g. backup the rotated file to backblaze.