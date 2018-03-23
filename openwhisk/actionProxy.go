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
)

// theServer is the current server
var theServer http.Server

// theChannel is the channel communicating with the action
var theExecutor *Executor

// StartLatestAction tries to start
// the more recently uplodaded
// action if valid, otherwise remove it
// and fallback to the previous, if any
func StartLatestAction() error {

	// find the action if any
	highestDir := highestDir("./action")
	if highestDir == 0 {
		log.Println("no action found")
		theExecutor = nil
		return fmt.Errorf("no valid actions available")
	}

	// save the current executor
	curExecutor := theExecutor

	// try to launch the action
	executable := fmt.Sprintf("./action/%d/exec", highestDir)
	newExecutor := NewExecutor(executable)
	log.Printf("starting %s", executable)
	err := newExecutor.Start()
	if err == nil {
		theExecutor = newExecutor
		if curExecutor != nil {
			log.Println("stopping old executor")
			curExecutor.Stop()
		}
		return nil
	}

	// cannot start, removing the action
	// and leaving the current executor running

	exeDir := fmt.Sprintf("./action/%d/", highestDir)
	log.Printf("removing the failed action in %s", exeDir)
	os.RemoveAll(exeDir)
	return err
}

// Start creates a proxy to execute actions
func Start() {
	// handle initialization
	http.HandleFunc("/init", initHandler)
	// handle execution
	http.HandleFunc("/run", runHandler)

	// start
	log.Println("Started!")
	theServer.Addr = ":8080"
	log.Fatal(theServer.ListenAndServe())
}
