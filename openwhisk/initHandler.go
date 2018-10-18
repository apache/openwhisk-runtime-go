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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type initBodyRequest struct {
	Code   string `json:"code,omitempty"`
	Binary bool   `json:"binary,omitempty"`
	Main   string `json:"main,omitempty"`
}

type initRequest struct {
	Value initBodyRequest `json:"value,omitempty"`
}

func sendOK(w http.ResponseWriter) {
	// answer OK
	w.Header().Set("Content-Type", "application/json")
	buf := []byte("{\"ok\":true}\n")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(buf)))
	w.Write(buf)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func (ap *ActionProxy) initHandler(w http.ResponseWriter, r *http.Request) {

	if ap.initialized {
		msg := "Cannot initialize the action more than once."
		sendError(w, http.StatusForbidden, msg)
		log.Println(msg)
		return
	}

	// read body of the request
	if ap.compiler != "" {
		Debug("compiler: " + ap.compiler)
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	// decode request parameters
	if len(body) < 1000 {
		Debug("init: decoding %s\n", string(body))
	}

	var request initRequest
	err = json.Unmarshal(body, &request)

	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error unmarshaling request: %v", err))
		return
	}

	// request with empty code - stop any executor but return ok
	if request.Value.Code == "" {
		sendError(w, http.StatusForbidden, "Missing main/no code to execute.")
		return
	}

	main := request.Value.Main
	if main == "" {
		main = "main"
	}

	// extract code eventually decoding it
	var buf []byte
	if request.Value.Binary {
		Debug("it is binary code")
		buf, err = base64.StdEncoding.DecodeString(request.Value.Code)
		if err != nil {
			sendError(w, http.StatusBadRequest, "cannot decode the request: "+err.Error())
			return
		}
	} else {
		Debug("it is source code")
		buf = []byte(request.Value.Code)
	}

	// if a compiler is defined try to compile
	_, err = ap.ExtractAndCompile(&buf, main)
	if err != nil {
		sendError(w, http.StatusBadGateway, err.Error())
		return
	}

	// start an action
	err = ap.StartLatestAction()
	if err != nil {
		sendError(w, http.StatusBadRequest, "cannot start action: "+err.Error())
		return
	}
	ap.initialized = true
	sendOK(w)
}

// ExtractAndCompile decode the buffer and if a compiler is defined, compile it also
func (ap *ActionProxy) ExtractAndCompile(buf *[]byte, main string) (string, error) {
	// extract in "bin" or in "src" if the runtime can compile
	suffix := "bin"
	if ap.compiler != "" {
		suffix = "src"
	}

	// extract action
	file, err := ap.ExtractAction(buf, suffix)
	if err != nil {
		return "", err
	}
	if file == "" {
		return "", fmt.Errorf("empty filename")
	}

	// some path surgery
	dir := filepath.Dir(file)
	parent := filepath.Dir(dir)
	srcDir := filepath.Join(parent, "src")
	binDir := filepath.Join(parent, "bin")
	binFile := filepath.Join(binDir, "exec")
	os.Mkdir(binDir, 0755)

	// if the file is already compiled just move it from src to bin
	if isCompiled(file) {
		os.Rename(file, binFile)
		return binFile, nil
	}

	// no compiler, move it anyway
	if ap.compiler == "" {
		os.Rename(file, binFile)
		return binFile, nil
	}

	// ok let's try to compile
	Debug("compiling: %s main: %s", file, main)
	err = ap.CompileAction(main, srcDir, binDir)
	if err != nil {
		return "", err
	}
	if !isCompiled(binFile) {
		return "", fmt.Errorf("cannot compile")
	}
	return binFile, nil
}
