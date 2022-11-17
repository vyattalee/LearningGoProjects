package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func errPanic2(_ http.ResponseWriter,
	_ *http.Request) error {
	panic(123)
}

type testingUserError string

func (e testingUserError) Error() string {
	return e.Message()
}

func (e testingUserError) Message() string {
	return string(e)
}

func errUserError(_ http.ResponseWriter,
	_ *http.Request) error {
	return testingUserError("user error")
}

func errNotFound(_ http.ResponseWriter,
	_ *http.Request) error {
	return os.ErrNotExist
}

func errNoPermission(_ http.ResponseWriter,
	_ *http.Request) error {
	return os.ErrPermission
}

func errUnknown(_ http.ResponseWriter,
	_ *http.Request) error {
	return errors.New("unknown error")
}

func noError(writer http.ResponseWriter,
	_ *http.Request) error {
	fmt.Fprintln(writer, "no error")
	return nil
}

var tests2 = []struct {
	h       appHandler
	code    int
	message string
}{
	{errPanic2, 500, "Internal Server Error"},
	{errUserError, 400, "user error"},
	{errNotFound, 404, "Not Found"},
	{errNoPermission, 403, "Forbidden"},
	{errUnknown, 500, "Internal Server Error"},
	{noError, 200, "no error"},
}

func TestErrWarpper2(t *testing.T) {
	for _, tt := range tests2 {
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

func TestErrWarpperInserver(t *testing.T) {
	for _, tt := range tests2 {
		f := errWarpper(tt.h)
		server := httptest.NewServer(http.HandlerFunc(f))
		resp, _ := http.Get(server.URL)
		b, _ := ioutil.ReadAll(resp.Body)
		body := strings.Trim(string(b), "\n") // web默认返回有个换行符，此处需要去掉才能正确匹配
		if resp.StatusCode != tt.code || body != tt.message {
			t.Errorf("expect (%d, %s); got (%d, %s)", tt.code, tt.message, resp.StatusCode, body)
		}
	}
}
