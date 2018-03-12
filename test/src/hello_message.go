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
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func hello(arg string) string {
	var obj map[string]interface{}
	json.Unmarshal([]byte(arg), &obj)
	name, ok := obj["name"].(string)
	if !ok {
		name = "Stranger"
	}
	log.Printf("name=%s\n", name)
	msg := map[string]string{"message": ("Hello, " + name + "!")}
	res, _ := json.Marshal(msg)
	return string(res)
}

func main() {
	log.SetPrefix("hello_message: ")
	log.SetFlags(0)
	// native actions receive one argument, the JSON object as a string
	if len(os.Args) > 1 {
		fmt.Print(hello(os.Args[1]))
		return
	}
	// read loop
	reader := bufio.NewReader(os.Stdin)
	for {
		event, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Println(hello(event))
	}
}
