package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func errPanic(_ http.ResponseWriter,
	_ *http.Request) error {
	panic(123)
}

var tests = []struct {
	h       appHandler
	code    int
	message string
}{
	{errPanic, 500, "Internal Server Error"},
}

func TestErrWarpper(t *testing.T) {
	for _, tt := range tests {
		f := errWarpper(tt.h)
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "http://www.imooc.com", nil)
		f(response, request)
		b, _ := ioutil.ReadAll(response.Body)
		body := strings.Trim(string(b), "\n") // web默认返回有个换行符，此处需要去掉才能正确匹配
		if response.Code != tt.code || body != tt.message {
			t.Errorf("expect (%d, %s); got (%d, %s)", tt.code, tt.message, response.Code, body)
		}
	}
}
