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
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/h2non/filetype"
	"github.com/stretchr/testify/assert"
)

func sys(cli string) {
	cmd := exec.Command(cli)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(err)
	} else {
		fmt.Print(string(out))
	}
}

func TestExtractActionTest_exec(t *testing.T) {
	sys("_test/build.sh")
	// cleanup
	assert.Nil(t, os.RemoveAll("./action"))
	file, _ := ioutil.ReadFile("_test/exec")
	extractAction(&file, false)
	assert.Nil(t, exists("./action", "exec"))
}

func exists(dir, filename string) error {
	path := fmt.Sprintf("%s/%d/%s", dir, highestDir(dir), filename)
	_, err := os.Stat(path)
	return err
}

func detect(dir, filename string) string {
	path := fmt.Sprintf("%s/%d/%s", dir, highestDir(dir), filename)
	file, _ := ioutil.ReadFile(path)
	kind, _ := filetype.Match(file)
	return kind.Extension
}
func TestExtractActionTest_exe(t *testing.T) {
	sys("_test/build.sh")
	// cleanup
	assert.Nil(t, os.RemoveAll("./action"))
	// match  exe
	file, _ := ioutil.ReadFile("_test/exec")
	extractAction(&file, false)
	assert.Equal(t, detect("./action", "exec"), "elf")
}

func TestExtractActionTest_zip(t *testing.T) {
	sys("_test/build.sh")
	// cleanup
	assert.Nil(t, os.RemoveAll("./action"))
	// match  exe
	file, _ := ioutil.ReadFile("_test/exec.zip")
	extractAction(&file, false)
	assert.Equal(t, detect("./action", "exec"), "elf")
	assert.Nil(t, exists("./action", "etc"))
	assert.Nil(t, exists("./action", "dir/etc"))
}

func TestExtractAction_script(t *testing.T) {
	buf := []byte("#!/bin/sh\necho ok")
	assert.NotNil(t, extractAction(&buf, false))
	assert.Nil(t, extractAction(&buf, true))
}

func TestHighestDir(t *testing.T) {
	assert.Equal(t, highestDir("./_test"), 0)
	assert.Equal(t, highestDir("./_test/first"), 3)
	assert.Equal(t, highestDir("./_test/second"), 17)
}
