package filelisting

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var prefix = "/list/"

type userError string

func (e userError) Error() string {
	return e.Message()
}

func (e userError) Message() string {
	return string(e)
}

func Handler(w http.ResponseWriter, r *http.Request) error {
	if strings.Index(r.URL.Path, prefix) != 0 {
		return userError("path must start with " + prefix)
	}

	path := r.URL.Path[len(prefix):] // /list/readme.txt

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	all, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	w.Write(all)

	return nil
}
