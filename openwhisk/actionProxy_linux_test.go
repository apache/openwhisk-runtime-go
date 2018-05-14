// +build linux

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

/*
This test depends on the fact the proxy can detect the termination
of an executable that terminates before or after reading the output.
On OSX command termination is not detected until some input is read.
*/

package openwhisk

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartLatestAction(t *testing.T) {

	// cleanup
	os.RemoveAll("./action")
	logf, _ := ioutil.TempFile("/tmp", "log")
	ap := NewActionProxy("./action", "", logf)

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
