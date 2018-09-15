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
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/apache/incubator-openwhisk-runtime-go/openwhisk"
)

// flag to enable debug
var debug = flag.Bool("debug", false, "enable debug output")

// flag to require on-the-fly compilation
var compile = flag.String("compile", "", "compile, reading in standard input the specified function, and producing the result in stdout")

// use the runtime as a compiler "on-the-fly"
func extractAndCompile(ap *openwhisk.ActionProxy) {
	// read the std in
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	// extract and compile it
	file, err := ap.ExtractAndCompile(&in, *compile)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	// read the file and write it to stdout
	out, err := ioutil.ReadFile(file)
	_, err = os.Stdout.Write(out)
	if err != nil {
		log.Printf(err.Error())
		return
	}
}

func main() {
	flag.Parse()

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

	// start the balls rolling
	openwhisk.Debug("OpenWhisk Go Proxy: starting")
	ap.Start(8080)

}
