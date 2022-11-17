package main

import (
	"GoProgressPerformance/http_performance/filelisting"
	"net/http"
	"os"

	"github.com/gpmgo/gopm/modules/log"
)

type appHandler func(writer http.ResponseWriter, request *http.Request) error

func errWarpper(handler appHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Info("Panic: %v", r)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		err := handler(writer, request)
		if err != nil {
			log.Warn("error occurred handling request: %s", err.Error())
			if userErr, ok := err.(userError); ok {
				http.Error(writer, userErr.Message(), http.StatusBadRequest)
				return
			}
			code := http.StatusOK
			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
				//http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}
			http.Error(writer, http.StatusText(code), code)
		}
	}
}

// 自定义用户 error接口
type userError interface {
	error
	Message() string
}

func main() {
	http.HandleFunc("/", errWarpper(filelisting.HandleFileList))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
