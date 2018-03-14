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
	"log"
	"net/http"
	"os"
	"os/exec"
)

// theServer is the current server
var theServer http.Server

// theChannel is the channel communicating with the action
var theChannel chan string

// theChannel is the channel communicating to flush the logs
var theLogger chan bool

func stopAction() {
	// terminate current action
	if theChannel != nil {
		log.Println("terminating old action")
		theChannel <- ""
		theChannel = nil
	}
	// terminate the logger
	if theLogger != nil {
		theLogger <- false
		theLogger = nil
	}
}

func reStartAction() error {
	// stop action if any
	stopAction()

	// find the action if any
	highestDir := highestDir("./action")
	if highestDir == 0 {
		log.Println("no action dir")
		theChannel = nil
		theLogger = nil
		return fmt.Errorf("no valid actions available")
	}

	// try to launch the action
	executable := fmt.Sprintf("./action/%d/exec", highestDir)
	_, err := exec.LookPath(executable)
	// try to start the action
	if err == nil {
		log.Printf("starting %s", executable)
		ch, chl := StartService(executable)
		if ch != nil {
			theChannel = ch
			theLogger = chl
			return nil
		}
	}

	// cannot start, removing the action and retry
	exeDir := fmt.Sprintf("./action/%d/", highestDir)
	os.RemoveAll(exeDir)
	reStartAction()
	return fmt.Errorf("sent invalid action")
}

// Start creates a proxy to execute actions
func Start() {
	// handle initialization
	http.HandleFunc("/init", initHandler)
	// handle execution
	http.HandleFunc("/run", runHandler)
	// start action if there
	reStartAction()
	// start
	log.Println("Started!")
	theServer.Addr = ":8080"
	log.Fatal(theServer.ListenAndServe())
}
