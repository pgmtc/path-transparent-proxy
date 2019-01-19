package httptesting

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func StartHttpServer(addr string, replyMessage string) {
	server := http.NewServeMux()
	server.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		message := replyMessage
		message += ":" + request.RequestURI
		writer.Write([]byte(message))
	})
	http.ListenAndServe(addr, server)
}

func AssertResponse(t *testing.T, message string, expectedResult string, url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("%s FAIL: %s", message, err.Error())
		return false
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("%s FAIL: %s", message, err.Error())
		return false
	}

	result := string(body)
	result = strings.Trim(result, " \n\r") // Trim

	if expectedResult != result {
		t.Errorf("%s: expected not equal to actual.\nExp:%s\nAct:%s\n", message, expectedResult, result)
		return false
	}
	return true
}