package filelisting

import (
	"io/ioutil"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path[len("/list/"):] // /list/readme.txt

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
