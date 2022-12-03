package main

import (
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/list/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/list/"):] // /list/readme.txt

		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		all, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		w.Write(all)
	})

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
