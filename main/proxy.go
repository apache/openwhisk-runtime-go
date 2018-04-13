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

	"github.com/sciabarracom/incubator-openwhisk-runtime-go/openwhisk"
)

// disable stderr except when debugging
var debug = flag.Bool("debug", false, "enable debug output")

var compiler = flag.String("compiler", os.Getenv("COMPILER"), "define the compiler on the command line")

func main() {
	flag.Parse()

	if !*debug {
		// hide log unless you are debugging
		log.SetOutput(ioutil.Discard)
	}
	// start the balls rolling
	ap := openwhisk.NewActionProxy("./action", *compiler, os.Stdout, !*debug)
	ap.Start(8080)
}
