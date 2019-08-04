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
	"fmt"
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

func main() {
	flag.Parse()

	// show version number
	if *version {
		fmt.Printf("OpenWhisk ActionLoop Proxy v%s\n", openwhisk.Version)
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
		ap.ExtractAndCompileIO(os.Stdin, os.Stdout, *compile)
		return
	}

	// start the balls rolling
	openwhisk.Debug("OpenWhisk ActionLoop Proxy %s: starting", openwhisk.Version)
	ap.Start(8080)

}
