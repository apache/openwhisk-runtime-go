/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	// {"error":"Missing main/no code to execute."}
	// {"error":"Error unmarshaling request: invalid character 'X' looking for beginning of value"}
	// {"error":"no action defined yet"}
	// {"error":"no action defined yet"}
}

func TestStartLatestAction_emit1(t *testing.T) {
	os.RemoveAll("./action/t2")
	logf, _ := ioutil.TempFile("/tmp", "log")
	ap := NewActionProxy("./action/t2", "", logf, logf)
	// start the action that emits 1
	buf := []byte("#!/bin/sh\nwhile read a; do echo 1 >&3 ; done\n")
	ap.ExtractAction(&buf, "bin")
	ap.StartLatestAction()
	res, _ := ap.theExecutor.Interact([]byte("x"))
	assert.Equal(t, res, []byte("1\n"))
	ap.theExecutor.Stop()
}

func TestStartLatestAction_terminate(t *testing.T) {
	os.RemoveAll("./action/t3")
	logf, _ := ioutil.TempFile("/tmp", "log")
	ap := NewActionProxy("./action/t3", "", logf, logf)
	// now start an action that terminate immediately
	buf := []byte("#!/bin/sh\ntrue\n")
	ap.ExtractAction(&buf, "bin")
	ap.StartLatestAction()
	assert.Nil(t, ap.theExecutor)
}

func TestStartLatestAction_emit2(t *testing.T) {
	os.RemoveAll("./action/t4")
	logf, _ := ioutil.TempFile("/tmp", "log")
	ap := NewActionProxy("./action/t4", "", logf, logf)
	// start the action that emits 2
	buf := []byte("#!/bin/sh\nwhile read a; do echo 2 >&3 ; done\n")
	ap.ExtractAction(&buf, "bin")
	ap.StartLatestAction()
	res, _ := ap.theExecutor.Interact([]byte("z"))
	assert.Equal(t, res, []byte("2\n"))
	/**/
	ap.theExecutor.Stop()
}
