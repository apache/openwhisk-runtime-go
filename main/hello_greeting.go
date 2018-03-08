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
	"fmt"
	"os"

	"github.com/sciabarracom/openwhisk-runtime-go/hello"
)

func main() {
	// handle command line argument
	if len(os.Args) > 1 {
		result, err := hello.Hello([]byte(os.Args[1]))
		if err == nil {
			fmt.Println(string(result))
			return
		}
		fmt.Printf("{ error: %q}\n", err.Error())
		return
	}
	// read loop
	fmt.Println(`{"openwhisk":1}`)
	reader := bufio.NewReader(os.Stdin)
	for {
		event, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		result, err := hello.Hello(event)
		if err != nil {
			fmt.Printf("{ error: %q}\n", err.Error())
			continue
		}
		fmt.Println(string(result))
	}
}
