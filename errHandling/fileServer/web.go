package main

import (
	filelisting "g-base/errHandling/fileServer/fileListing"
	"log"
	"net/http"
	"os"
)

type appHandler func(w http.ResponseWriter, r *http.Request) error

func errWrapper(handler appHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			log.Println("Error handling request:", err.Error())
			code := http.StatusOK
			switch {
			case os.IsNotExist(err): // 文件不存在的错误
				code = http.StatusNotFound
			case os.IsPermission(err): // 没权限读
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}
			http.Error(w, http.StatusText(code), code)
		}
	}
}

func main() {
	http.HandleFunc("/list/", errWrapper(filelisting.Handler))

	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
