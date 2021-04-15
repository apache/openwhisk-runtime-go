
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
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
)

type Action func(event map[string]interface{}) map[string]interface{}


/*
	The Main function Simuates the action loop
	to pass in parameers add the --param flag and enter a json object with your parameters
	e.g. '{"name":"Test"}'
*/
func main() {
	var action Action
	// Action to call is set here
	action = Main

	// Get Parametes using --param flag
	paramsString := flag.String("param", "{}", "Parameters in JSON format")
	flag.Parse()

	// Unmarshal JSON
	var params map[string]interface{}
	err := json.Unmarshal([]byte(*paramsString), &params)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Run Action with Parameters
	result := action(params)

	// encode answer to JSON
	resultJSON, err := json.Marshal(&result)
	if err != nil {
		fmt.Println(err.Error())
	}
	resultString := string(bytes.Replace(resultJSON, []byte("\n"), []byte(""), -1))
	fmt.Println(resultString)


}

// Main function for the action
func Main(obj map[string]interface{}) map[string]interface{} {
	name, ok := obj["name"].(string)
	if !ok {
		name = "world"
	}
	fmt.Printf("name=%s\n", name)
	msg := make(map[string]interface{})
	msg["standalone"] = "Hello, " + name + "!"
	return msg
}
