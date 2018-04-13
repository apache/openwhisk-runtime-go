package openwhisk

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Example_startTestServer() {
	ts, cur, log := startTestServer("")
	res, _, _ := doPost(ts.URL+"/init", "{}")
	fmt.Print(res)
	res, _, _ = doPost(ts.URL+"/init", "XXX")
	fmt.Print(res)
	res, _, _ = doPost(ts.URL+"/run", "{}")
	fmt.Print(res)
	res, _, _ = doPost(ts.URL+"/run", "XXX")
	fmt.Print(res)
	stopTestServer(ts, cur, log)
	// Output:
	// {"ok":true}
	// {"error":"Error unmarshaling request: invalid character 'X' looking for beginning of value"}
	// {"error":"no action defined yet"}
	// {"error":"Error unmarshaling request: invalid character 'X' looking for beginning of value"}
}

func TestStartLatestAction(t *testing.T) {

	// cleanup
	os.RemoveAll("./action")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action", "", log, true)

	// start an action that terminate immediately
	buf := []byte("#!/bin/sh\ntrue\n")
	ap.ExtractAction(&buf, "main")
	ap.StartLatestAction("main")
	assert.Nil(t, ap.theExecutor)

	// start the action that emits 1
	buf = []byte("#!/bin/sh\nwhile read a; do echo 1 >&3 ; done\n")
	ap.ExtractAction(&buf, "main")
	ap.StartLatestAction("main")
	ap.theExecutor.io <- "x"
	assert.Equal(t, <-ap.theExecutor.io, "1")

	// now start an action that terminate immediately
	buf = []byte("#!/bin/sh\ntrue\n")
	ap.ExtractAction(&buf, "main")
	ap.StartLatestAction("main")
	ap.theExecutor.io <- "y"
	assert.Equal(t, <-ap.theExecutor.io, "1")

	// start the action that emits 2
	buf = []byte("#!/bin/sh\nwhile read a; do echo 2 >&3 ; done\n")
	ap.ExtractAction(&buf, "main")
	ap.StartLatestAction("main")
	ap.theExecutor.io <- "z"
	assert.Equal(t, <-ap.theExecutor.io, "2")
	/**/
	ap.theExecutor.Stop()
}
