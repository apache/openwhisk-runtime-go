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

// ActionProxy is the container of the data specific to a server
type ActionProxy struct {
	// current directory
	baseDir string

	// Compiler is the script to use to compile your code when action are source code
	compiler string

	// index current dir
	currentDir int

	// theChannel is the channel communicating with the action
	theExecutor *Executor

	// log file
	logFile *os.File

	// debug
	Debug bool
	Trace bool
}

// NewActionProxy creates a new action proxy that can handle http requests
func NewActionProxy(baseDir string, compiler string, logFile *os.File) *ActionProxy {
	os.Mkdir(baseDir, 0755)
	return &ActionProxy{
		baseDir,
		compiler,
		highestDir(baseDir),
		nil,
		logFile,
		false,
		false,
	}
}

// StartLatestAction tries to start
// the more recently uplodaded
// action if valid, otherwise remove it
// and fallback to the previous, if any
func (ap *ActionProxy) StartLatestAction(main string) error {

	// find the action if any
	highestDir := highestDir(ap.baseDir)
	if highestDir == 0 {
		log.Println("no action found")
		ap.theExecutor = nil
		return fmt.Errorf("no valid actions available")
	}

	// save the current executor
	curExecutor := ap.theExecutor

	// try to launch the action
	executable := fmt.Sprintf("%s/%d/%s", ap.baseDir, highestDir, main)
	os.Chmod(executable, 0755)
	newExecutor := NewExecutor(ap.logFile, executable)
	log.Printf("starting %s", executable)
	err := newExecutor.Start()
	if err == nil {
		ap.theExecutor = newExecutor
		if curExecutor != nil {
			log.Println("stopping old executor")
			curExecutor.Stop()
		}
		return nil
	}

	// cannot start, removing the action
	// and leaving the current executor running
	if !ap.Debug {
		exeDir := fmt.Sprintf("./action/%d/", highestDir)
		log.Printf("removing the failed action in %s", exeDir)
		os.RemoveAll(exeDir)
	}

	return err
}

func (ap *ActionProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(r.URL.Path)
	switch r.URL.Path {
	case "/init":
		ap.initHandler(w, r)
	case "/run":
		ap.runHandler(w, r)
	}
}

// Start creates a proxy to execute actions
func (ap *ActionProxy) Start(port int) {

	// listen and start
	//http.HandleFunc("/init", func(w http.ResponseWriter, r *http.Request) { ap.initHandler(w, r) })
	//http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) { ap.runHandler(w, r) })
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), ap))
}
