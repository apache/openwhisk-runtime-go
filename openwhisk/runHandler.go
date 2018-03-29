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
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Params are the parameteres sent to the action
type Params struct {
	Value json.RawMessage `json:"value"`
}

// ErrResponse is the response when there are errors
type ErrResponse struct {
	Error string `json:"error"`
}

func sendError(w http.ResponseWriter, code int, cause string) {
	errResponse := ErrResponse{Error: cause}
	b, err := json.Marshal(errResponse)
	if err != nil {
		b = []byte("error marshalling error response")
		fmt.Println(b, err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
	w.Write([]byte("\n"))
}

func runHandler(w http.ResponseWriter, r *http.Request) {

	// parse the request
	params := Params{}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error reading request body: %v", err))
		return
	}

	// decode request parameters
	err = json.Unmarshal(body, &params)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error unmarshaling request: %v", err))
		return
	}

	// check if you have an action
	if theExecutor == nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("no action defined yet"))
		return
	}

	// execute the action
	// and check for early termination
	theExecutor.io <- string(params.Value)
	var response string
	var exited bool
	select {
	case response = <-theExecutor.io:
		exited = false
	case err = <-theExecutor.exit:
		exited = true
	}

	// check for early termination
	if exited {
		theExecutor = nil
		sendError(w, http.StatusBadRequest, fmt.Sprintf("command exited"))
		return
	}

	// flush the logs sending the activation message at the end
	theExecutor.log <- true

	// check response
	if response == "" {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	// return the response
	if !strings.HasSuffix(response, "\n") {
		response = response + "\n"
	}
	log.Print(response)
	w.Header().Set("Content-Type", "application/json")
	numBytesWritten, err := w.Write([]byte(response))

	// flush output
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// diagnostic when writing problems
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error writing response: %v", err))
		return
	}
	if numBytesWritten != len(response) {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Only wrote %d of %d bytes to response", numBytesWritten, len(response)))
		return
	}
}
