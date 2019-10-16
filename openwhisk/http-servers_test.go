package openwhisk

import (
	"fmt"
	"net/http"
)

func ExampleHello() {
	ts, cur, log := startTestServer("")
	res, _ := http.NewRequest("GET", "/hello", nil)
	fmt.Print(res)
	stopTestServer(ts, cur, log)
	// Output: hello
}
