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
	"reflect"
	"strings"
)

func main() {
	// debugging
	var debug = os.Getenv("OW_DEBUG") != ""

	if debug {
		filename := os.Getenv("OW_DEBUG")
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil {
			log.SetOutput(f)
			defer f.Close()
		}
		log.Printf("ACTION ENV: %v", os.Environ())
	}

	resultKind := reflect.TypeOf(Main).Out(0).Kind()
	if resultKind != reflect.Map && resultKind != reflect.Slice && resultKind != reflect.Array {
		fmt.Println("Support map and slice and array only")
		os.Exit(1)
	}

	// input
	out := os.NewFile(3, "pipe")
	defer out.Close()
	reader := bufio.NewReader(os.Stdin)

	// read-eval-print loop
	if debug {
		log.Println("started")
	}
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
		isJsonObjectParam := true
		var payloadForJsonObject map[string]interface{}
		var payloadForJsonArray []interface{}
		if value, ok := input["value"].(map[string]interface{}); ok {
			payloadForJsonObject = value
		} else {
			if value, ok := input["value"].([]interface{}); ok {
				payloadForJsonArray = value
				isJsonObjectParam = false
			}
		}
		// process the request
		var result interface{}
		funcMain := reflect.ValueOf(Main)
		if isJsonObjectParam {
			param := []reflect.Value{reflect.ValueOf(payloadForJsonObject)}
			reflectResult := funcMain.Call(param)
			result = reflectResult[0].Interface()
		} else {
			param := []reflect.Value{reflect.ValueOf(payloadForJsonArray)}
			reflectResult := funcMain.Call(param)
			result = reflectResult[0].Interface()
		}
		// encode the answer
		output, err := json.Marshal(&result)
		if err != nil {
			log.Println(err.Error())
			fmt.Fprintf(out, "{ error: %q}\n", err.Error())
			continue
		}
		output = bytes.Replace(output, []byte("\n"), []byte(""), -1)
		if debug {
			log.Printf("'<<<%s'<<<", output)
		}
		fmt.Fprintf(out, "%s\n", output)
	}
}
