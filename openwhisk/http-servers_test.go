package openwhisk

import (
	"fmt"
)

func ExampleHello() {
	ts, cur, log := startTestServer("")
	res, _, _ := doGet(ts.URL + "/hello")
	fmt.Println(res)
	stopTestServer(ts, cur, log)
	// Output:
	// hello
}
