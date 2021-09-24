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
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// OwExecutionEnv is the execution environment set at compile time
var OwExecutionEnv = ""

func main() {
	// check if the execution environment is correct
	if OwExecutionEnv != "" && OwExecutionEnv != os.Getenv("__OW_EXECUTION_ENV") {
		fmt.Println("Execution Environment Mismatch")
		fmt.Println("Expected: ", OwExecutionEnv)
		fmt.Println("Actual: ", os.Getenv("__OW_EXECUTION_ENV"))
		os.Exit(1)
	}

	// debugging
	var debug = os.Getenv("OW_DEBUG") != ""
	if debug {
		f, err := os.OpenFile("/tmp/action.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(f)
		}
		log.Printf("Environment: %v", os.Environ())
	}

	// assign the main function
	type Action func(event map[string]interface{}) map[string]interface{}
	var action Action
	action = Main

	// input
	out := os.NewFile(3, "pipe")
	defer out.Close()
	reader := bufio.NewReader(os.Stdin)

	// acknowledgement of started action
	fmt.Fprintf(out, `{ "ok": true}%s`, "\n")
	if debug {
		log.Println("action started")
	}

	// read-eval-print loop
	for {
		// read one line
		inbuf, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}
		if debug {
			log.Printf(">>>'%s'>>>", inbuf)
		}
		// parse one line
		var input map[string]interface{}
		err = json.Unmarshal(inbuf, &input)
		if err != nil {
			log.Println(err.Error())
			fmt.Fprintf(out, "{ error: %q}\n", err.Error())
			continue
		}
		if debug {
			log.Printf("%v\n", input)
		}
		// set environment variables
		for k, v := range input {
			if k == "value" {
				continue
			}
			if s, ok := v.(string); ok {
				os.Setenv("__OW_"+strings.ToUpper(k), s)
			}
		}
		// get payload if not empty
		var payload map[string]interface{}
		if value, ok := input["value"].(map[string]interface{}); ok {
			payload = value
		}
		// process the request
		result := action(payload)
		// encode the answer
		output, err := json.Marshal(&result)
		if err != nil {
			log.Println(err.Error())
			fmt.Fprintf(out, "{ error: %q}\n", err.Error())
			continue
		}
		output = bytes.Replace(output, []byte("\n"), []byte(""), -1)
		if debug {
			log.Printf("<<<'%s'<<<", output)
		}
		fmt.Fprintf(out, "%s\n", output)
	}
}
