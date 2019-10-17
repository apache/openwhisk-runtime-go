package openwhisk

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleHealthCheckHandler() {

	ts, cur, log := startTestServer("")

	req, err := http.NewRequest("GET", "hello", nil)
	if err != nil {
		fmt.Println(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(startTestServer)

	handler.ServeHTTP(rr, req)

	// Output: hello

	stopTestServer(ts, cur, log)
}
