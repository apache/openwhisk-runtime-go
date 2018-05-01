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

package action

import (
	"encoding/json"
	"fmt"
	//"log"
)

// Hello receive an event in format
// { "name": "Mike"}
// and returns a greeting in format
// { "greetings": "Hello, Mike"}
func Hello(event json.RawMessage) (json.RawMessage, error) {
	// input and output
	var input struct {
		Name string
	}
	var output struct {
		Greetings string `json:"greetings"`
	}
	// read the input event
	json.Unmarshal(event, &input)
	if input.Name != "" {
		// handle the event
		output.Greetings = "Hello, " + input.Name
		fmt.Println(output.Greetings)
		return json.Marshal(output)
	}
	return nil, fmt.Errorf("no name specified")
}
