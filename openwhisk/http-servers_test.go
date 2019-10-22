package openwhisk

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func ExampleHello(w http.ResponseWriter, r *http.Request) {

	ts, cur, log := startTestServer("")

	res, err := http.Get("/hello")
	if err != nil {
		println(err)
	}
	defer res.Body.Close()

	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		println(err)
	}

	responseString := string(responseData)
	fmt.Fprint(w, responseString)

	println(res)

	// Output: hello

	stopTestServer(ts, cur, log)
}
