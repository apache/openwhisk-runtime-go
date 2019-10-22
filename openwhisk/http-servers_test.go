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

	responseData, err := ioutil.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		println(err)
	}

	responseString := string(responseData)
	fmt.Fprint(w, responseString)

	println(res)

	// Output: hello

	stopTestServer(ts, cur, log)
}
