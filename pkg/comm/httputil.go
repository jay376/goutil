package comm

import (
	"io"
	"net/http"
	"os"
)

// Download ...
func Download(url, f string) error {
	res, err := http.Get(url) //nolint
	if err == nil {
		defer res.Body.Close()
		tmpFile := f + ".tmp"
		os.Remove(tmpFile)
		var tmp *os.File
		if tmp, err = os.Create(tmpFile); err == nil {
			if _, err = io.Copy(tmp, res.Body); err == nil {
				tmp.Close()
				err = os.Rename(tmpFile, f)
			}
		}
	}
	return err
}
