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

	// read body of the request
	if ap.compiler != "" {
		log.Println("compiler: " + ap.compiler)
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	// decode request parameters
	if ap.Trace {
		log.Printf("init: decoding %s\n", string(body))
	}
	var request initRequest
	err = json.Unmarshal(body, &request)

	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error unmarshaling request: %v", err))
		return
	}

	// request with empty code - stop any executor but return ok
	if ap.Trace {
		log.Printf("request: %v\n", request)
	}
	if request.Value.Code == "" {
		if ap.theExecutor != nil {
			log.Printf("stop running action")
			ap.theExecutor.Stop()
			ap.theExecutor = nil
		}
		sendOK(w)
		return
	}

	main := request.Value.Main
	if main == "" {
		main = "main"
	}

	// extract code eventually decoding it
	var buf []byte
	if request.Value.Binary {
		log.Printf("binary")
		buf, err = base64.StdEncoding.DecodeString(request.Value.Code)
		if err != nil {
			sendError(w, http.StatusBadRequest, "cannot decode the request: "+err.Error())
			return
		}
	} else {
		log.Printf("plain text")
		buf = []byte(request.Value.Code)
	}

	// extract the action,
	file, err := ap.ExtractAction(&buf, main)
	if err != nil || file == "" {
		sendError(w, http.StatusBadRequest, "invalid action: "+err.Error())
		return
	}

	// compile it if a compiler is available
	if ap.compiler != "" && !isCompiled(file, main) {
		log.Printf("compiling: %s main: %s", file, main)
		err = ap.CompileAction(main, file, file)
		if err != nil {
			sendError(w, http.StatusBadRequest, "cannot compile action: "+err.Error())
			return
		}
	}

	// stop and start
	err = ap.StartLatestAction(main)
	if err != nil {
		sendError(w, http.StatusBadRequest, "cannot start action: "+err.Error())
		return
	}
	sendOK(w)
}
