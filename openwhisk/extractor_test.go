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

func TestExtractActionTest_exec(t *testing.T) {
	ap := NewActionProxy("./action/x1", "", os.Stdout, os.Stderr)
	// cleanup
	assert.Nil(t, os.RemoveAll("./action/x1"))
	file, _ := ioutil.ReadFile("_test/exec")
	ap.ExtractAction(&file, "bin")
	assert.Nil(t, exists("./action/x1", "bin/exec"))
}

func TestExtractActionTest_exe(t *testing.T) {
	ap := NewActionProxy("./action/x2", "", os.Stdout, os.Stderr)
	// cleanup
	assert.Nil(t, os.RemoveAll("./action/x2"))
	// match  exe
	file, _ := ioutil.ReadFile("_test/exec")
	ap.ExtractAction(&file, "bin")
	assert.Equal(t, detectExecutable("./action/x2", "bin/exec"), true)
}

func TestExtractActionTest_zip(t *testing.T) {
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/x3", "", log, log)
	// cleanup
	assert.Nil(t, os.RemoveAll("./action/x3"))
	// match  exe
	file, _ := ioutil.ReadFile("_test/exec.zip")
	ap.ExtractAction(&file, "bin")
	assert.Equal(t, detectExecutable("./action/x3", "bin/exec"), true)
	assert.Nil(t, exists("./action/x3", "bin/etc"))
	assert.Nil(t, exists("./action/x3", "bin/dir/etc"))
}

func TestExtractAction_script(t *testing.T) {
	log, _ := ioutil.TempFile("", "log")
	assert.Nil(t, os.RemoveAll("./action/x4"))
	ap := NewActionProxy("./action/x4", "", log, log)
	buf := []byte("#!/bin/sh\necho ok")
	_, err := ap.ExtractAction(&buf, "bin")
	//fmt.Print(err)
	assert.Nil(t, err)
}

func TestExtractAction_save_jar(t *testing.T) {
	os.Setenv("OW_SAVE_JAR", "exec.jar")
	log, _ := ioutil.TempFile("", "log")
	assert.Nil(t, os.RemoveAll("./action/x5"))
	ap := NewActionProxy("./action/x5", "", log, log)
	file, _ := ioutil.ReadFile("_test/sample.jar")
	_, err := ap.ExtractAction(&file, "bin")
	assert.Nil(t, exists("./action/x5", "bin/exec.jar"))
	assert.Nil(t, err)
	os.Setenv("OW_SAVE_JAR", "")
}

func TestExtractAction_extract_jar(t *testing.T) {
	os.Setenv("OW_SAVE_JAR", "")
	log, _ := ioutil.TempFile("", "log")
	assert.Nil(t, os.RemoveAll("./action/x6"))
	ap := NewActionProxy("./action/x6", "", log, log)
	file, _ := ioutil.ReadFile("_test/sample.jar")
	_, err := ap.ExtractAction(&file, "bin")
	assert.Nil(t, exists("./action/x6", "bin/META-INF/MANIFEST.MF"))
	assert.Nil(t, err)
}


func TestHighestDir(t *testing.T) {
	assert.Equal(t, highestDir("./_test"), 0)
	assert.Equal(t, highestDir("./_test/first"), 3)
	assert.Equal(t, highestDir("./_test/second"), 17)
}

// Issue #62 sample zip
func Example_badZip() {
	buf := []byte{
		0x50, 0x4b, 0x03, 0x04, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x81, 0xb0, 0x81, 0x4d, 0x2d, 0xf6,
		0xa5, 0x66, 0x48, 0x00, 0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x09, 0x00, 0x1c, 0x00, 0x69, 0x6e,
		0x64, 0x65, 0x78, 0x2e, 0x70, 0x68, 0x70, 0x55, 0x54, 0x09, 0x00, 0x03, 0x51, 0x05, 0x03, 0x5c,
		0x54, 0x05, 0x03, 0x5c, 0x75, 0x78, 0x0b, 0x00, 0x01, 0x04, 0xf5, 0x01, 0x00, 0x00, 0x04, 0x14,
		0x00, 0x00, 0x00, 0x3c, 0x3f, 0x70, 0x68, 0x70, 0x0a, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f,
		0x6e, 0x20, 0x6e, 0x69, 0x61, 0x6d, 0x28, 0x61, 0x72, 0x72, 0x61, 0x79, 0x20, 0x24, 0x61, 0x72,
		0x67, 0x73, 0x29, 0x20, 0x7b, 0x0a, 0x20, 0x20, 0x20, 0x20, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e,
		0x20, 0x5b, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x20, 0x3d, 0x3e, 0x20, 0x22, 0x69, 0x74, 0x20,
		0x77, 0x6f, 0x72, 0x6b, 0x73, 0x22, 0x5d, 0x3b, 0x0a, 0x7d, 0x0a, 0x50, 0x4b, 0x01, 0x02, 0x1e,
		0x03, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x81, 0xb0, 0x81, 0x4d, 0x2d, 0xf6, 0xa5, 0x66, 0x48,
		0x00, 0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x09, 0x00, 0x18, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0xb4, 0x81, 0x00, 0x00, 0x00, 0x00, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x70,
		0x68, 0x70, 0x55, 0x54, 0x05, 0x00, 0x03, 0x51, 0x05, 0x03, 0x5c, 0x75, 0x78, 0x0b, 0x00, 0x01,
		0x04, 0xf5, 0x01, 0x00, 0x00, 0x04, 0x14, 0x00, 0x00, 0x00, 0x50, 0x4b, 0x05, 0x06, 0x00, 0x00,
		0x00, 0x00, 0x01, 0x00, 0x01, 0x00, 0x4f, 0x00, 0x00, 0x00, 0x8b, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	fmt.Printf("%t", IsZip(buf))
	// Output:
	// true
}
