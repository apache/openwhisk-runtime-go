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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ActionProxy is the container of the data specific to a server
type ActionProxy struct {

	// is it initialized?
	initialized bool

	// current directory
	baseDir string

	// Compiler is the script to use to compile your code when action are source code
	compiler string

	// index current dir
	currentDir int

	// theChannel is the channel communicating with the action
	theExecutor *Executor

	// out and err files
	outFile *os.File
	errFile *os.File

	// environment
	env map[string]string
}

// NewActionProxy creates a new action proxy that can handle http requests
func NewActionProxy(baseDir string, compiler string, outFile *os.File, errFile *os.File) *ActionProxy {
	os.Mkdir(baseDir, 0755)
	return &ActionProxy{
		false,
		baseDir,
		compiler,
		highestDir(baseDir),
		nil,
		outFile,
		errFile,
		map[string]string{},
	}
}

//SetEnv sets the environment
func (ap *ActionProxy) SetEnv(env map[string]interface{}) {
	// Propagate proxy version
	ap.env["__OW_PROXY_VERSION"] = Version
	// propagate OW_EXECUTION_ENV as  __OW_EXECUTION_ENV
	ee := os.Getenv("OW_EXECUTION_ENV")
	if ee != "" {
		ap.env["__OW_EXECUTION_ENV"] = ee
	}
	// require an ack
	wa := os.Getenv("OW_WAIT_FOR_ACK")
	if wa != "" {
		ap.env["__OW_WAIT_FOR_ACK"] = wa
	}
	// propagate all the variables starting with "__OW_"
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, "__OW_") {
			res := strings.Split(v, "=")
			ap.env[res[0]] = res[1]
		}
	}
	// get other variables from the init payload
	for k, v := range env {
		s, ok := v.(string)
		if ok {
			ap.env[k] = s
			continue
		}
		buf, err := json.Marshal(v)
		if err == nil {
			ap.env[k] = string(buf)
		}
	}
	Debug("init env: %s", ap.env)
}

// StartLatestAction tries to start
// the more recently uploaded
// action if valid, otherwise remove it
// and fallback to the previous, if any
func (ap *ActionProxy) StartLatestAction() error {

	// find the action if any
	highestDir := highestDir(ap.baseDir)
	if highestDir == 0 {
		Debug("no action found")
		ap.theExecutor = nil
		return fmt.Errorf("no valid actions available")
	}

	// check version
	execEnv := os.Getenv("OW_EXECUTION_ENV")
	if execEnv != "" {
		execEnvFile := fmt.Sprintf("%s/%d/bin/exec.env", ap.baseDir, highestDir)
		execEnvData, err := ioutil.ReadFile(execEnvFile)
		if err != nil {
			return err
		}
		if strings.TrimSpace(string(execEnvData)) != execEnv {
			fmt.Printf("Expected exec.env should start with %s\nActual value: %s", execEnv, execEnvData)
			return fmt.Errorf("Execution environment version mismatch. See logs for details.")
		}
	}

	// save the current executor
	curExecutor := ap.theExecutor

	// try to launch the action
	executable := fmt.Sprintf("%s/%d/bin/exec", ap.baseDir, highestDir)
	os.Chmod(executable, 0755)
	newExecutor := NewExecutor(ap.outFile, ap.errFile, executable, ap.env)
	Debug("starting %s", executable)

	// start executor
	err := newExecutor.Start(os.Getenv("OW_WAIT_FOR_ACK") != "")
	if err == nil {
		ap.theExecutor = newExecutor
		if curExecutor != nil {
			Debug("stopping old executor")
			curExecutor.Stop()
		}
		return nil
	}

	// cannot start, removing the action
	// and leaving the current executor running
	if !Debugging {
		exeDir := fmt.Sprintf("./action/%d/", highestDir)
		Debug("removing the failed action in %s", exeDir)
		os.RemoveAll(exeDir)
	}
	return err
}

func (ap *ActionProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), ap))
}

// ExtractAndCompileIO read in input and write in output to use the runtime as a compiler "on-the-fly"
func (ap *ActionProxy) ExtractAndCompileIO(r io.Reader, w io.Writer, main string, env string) {

	// read the std input
	in, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	envMap := make(map[string]interface{})
	if env != "" {
	    json.Unmarshal([]byte(env), &envMap)
	}
    ap.SetEnv(envMap)

	// extract and compile it
	file, err := ap.ExtractAndCompile(&in, main)
	if err != nil {
		log.Fatal(err)
	}

	// zip the directory containing the file and write output
	zip, err := Zip(filepath.Dir(file))
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(zip)
	if err != nil {
		log.Fatal(err)
	}
}
