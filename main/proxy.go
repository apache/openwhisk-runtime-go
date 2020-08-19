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
package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/apache/openwhisk-runtime-go/openwhisk"
)

// flag to show version
var version = flag.Bool("version", false, "show version")

// flag to enable debug
var debug = flag.Bool("debug", false, "enable debug output")

// flag to require on-the-fly compilation
var compile = flag.String("compile", "", "compile, reading in standard input the specified function, and producing the result in stdout")

// fatal if error
func fatalIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// use the runtime as a compiler "on-the-fly"
func extractAndCompile(ap *openwhisk.ActionProxy) {

	// read the std input
	in, err := ioutil.ReadAll(os.Stdin)
	fatalIf(err)

	// extract and compile it
	file, err := ap.ExtractAndCompile(&in, *compile)
	fatalIf(err)

	// read the file, zip it and write it to stdout
	buf := new(bytes.Buffer)
	zwr := zip.NewWriter(buf)
	zf, err := zwr.Create("exec")
	fatalIf(err)
	filedata, err := ioutil.ReadFile(file)
	fatalIf(err)
	_, err = zf.Write(filedata)
	fatalIf(err)
	fatalIf(zwr.Flush())
	fatalIf(zwr.Close())
	_, err = os.Stdout.Write(buf.Bytes())
	fatalIf(err)
}

// auto init
func autoInit(ap *openwhisk.ActionProxy, initFile string, main string) {
	if initFile == "" {
		return
	}
	_, err := os.Stat(initFile)
	if os.IsNotExist(err) {
		openwhisk.Debug("not found %s: %s", initFile, err)
		return
	}
	buf, err := ioutil.ReadFile(initFile)
	if err != nil {
		openwhisk.Debug("cannot read: %s", err)
		return
	}

	// default main to "main"
	if main == "" {
		main = "main"
	}

	// if a compiler is defined try to compile
	_, err = ap.ExtractAndCompile(&buf, main)
	if err != nil {
		openwhisk.Debug("compile error: %s", err)
		return
	}

	// start an action
	err = ap.StartLatestAction()
	if err != nil {
		openwhisk.Debug("%s", err)
		return
	}
	openwhisk.Debug("action autoinitialized")
	ap.Initialized = true
}

func main() {
	flag.Parse()

	// show version number
	if *version {
		fmt.Println("OpenWhisk ActionLoop Proxy", openwhisk.Version)
		return
	}

	// debugging
	if *debug {
		// set debugging flag, propagated to the actions
		openwhisk.Debugging = true
		os.Setenv("OW_DEBUG", "1")
	}

	// create the action proxy
	ap := openwhisk.NewActionProxy("./action", os.Getenv("OW_COMPILER"), os.Stdout, os.Stderr)

	// compile on the fly upon request
	if *compile != "" {
		extractAndCompile(ap)
		return
	}

	// auto initialization if available
	autoInit(ap, os.Getenv("OW_AUTOINIT"), os.Getenv("OW_AUTOINIT_MAIN"))

	// start the balls rolling
	openwhisk.Debug("OpenWhisk ActionLoop Proxy %s: starting", openwhisk.Version)
	ap.Start(8080)

}
