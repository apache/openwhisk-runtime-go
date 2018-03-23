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
	"net/http"
)

type initRequest struct {
	Value struct {
		Code   string `json:",omitempty"`
		Binary bool   `json:",omitempty"`
	} `json:",omitempty"`
}

func sendOK(w http.ResponseWriter) {
	// answer OK
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", "12")
	w.Write([]byte("{\"ok\":true}\n"))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func initHandler(w http.ResponseWriter, r *http.Request) {

	// read body of the request
	// log.Println("init: reading")
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	// decode request parameters
	//log.Println("init: decoding")
	var request initRequest
	err = json.Unmarshal(body, &request)

	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error unmarshaling request: %v", err))
		return
	}

	if request.Value.Code == "" {
		sendOK(w)
		return
	}

	// check if it is a binary
	if request.Value.Binary {
		var decoded []byte
		decoded, err = base64.StdEncoding.DecodeString(request.Value.Code)
		if err != nil {
			sendError(w, http.StatusBadRequest, "cannot decode the request: "+err.Error())
			return
		}
		// extract the replacement, stopping and then starting the action
		err = extractAction(&decoded, false)
	} else {
		buf := []byte(request.Value.Code)
		err = extractAction(&buf, true)
	}
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid action: "+err.Error())
		return
	}

	// stop and start
	err = StartLatestAction()
	if err != nil {
		sendError(w, http.StatusBadRequest, "cannot start action: "+err.Error())
		return
	}
	sendOK(w)
}
