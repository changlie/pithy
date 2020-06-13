package pithy

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func localFile(path string) (string, error) {
	if path == "" || path == "/" {
		return "", errors.New("don't handle root path")
	}
	fp := filepath.Join(Pwd(), "static", path)
	_, err := os.Open(fp)
	return fp, err
}

func returnLocalFile(w http.ResponseWriter, fp string) {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	w.Write(data)
}
