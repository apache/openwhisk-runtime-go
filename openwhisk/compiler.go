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
	"runtime"

	"github.com/h2non/filetype"
)

// this is only to let test run on OSX
// it only recognizes OSX Mach 64 bit executable
// (magic number: facefeed + 64bit flag)
var mach64Type = filetype.NewType("mach", "darwin/mach")

func mach64Matcher(buf []byte) bool {
	return len(buf) > 4 && buf[0] == 0xcf && buf[1] == 0xfa && buf[2] == 0xed && buf[3] == 0xfe
}

// check if the file is already compiled
// if the file is a directoy look for a file with the given name
func isCompiled(fileOrDir string, name string) bool {
	fi, err := os.Stat(fileOrDir)
	if err != nil {
		log.Print(err)
		return false
	}
	file := fileOrDir
	if fi.IsDir() {
		file = fmt.Sprintf("%s/%s", fileOrDir, name)
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return false
	}
	// if a mac add a matched for mac
	if runtime.GOOS == "darwin" {
		filetype.AddMatcher(mach64Type, mach64Matcher)
	}

	kind, err := filetype.Match(buf)
	log.Printf("isCompiled: %s kind=%s", file, kind)
	if err != nil {
		log.Println(err)
		return false
	}
	if kind.Extension == "elf" {
		return true
	}
	if kind.Extension == "mach" {
		return true
	}
	return false
}

// CompileAction will compile an anction in source format invoking a compiler
func (ap *ActionProxy) CompileAction(main string, src string, out string) error {
	if ap.compiler == "" {
		return fmt.Errorf("No compiler defined")
	}
	log.Printf("compiling: compiler=%s main=%s src=%s out=%s", ap.compiler, main, src, out)
	var cmd *exec.Cmd
	if out == "" {
		cmd = exec.Command(ap.compiler, main, src, src)
	} else {
		cmd = exec.Command(ap.compiler, main, src, out)
	}
	res, err := cmd.CombinedOutput()
	log.Print(string(res))
	return err
}
