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
	"fmt"
	"io/ioutil"
)

func ExampleNewExecutor_failed() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, "true")
	err := proc.Start()
	fmt.Println(err)
	proc.Stop()
	proc = NewExecutor(log, "/bin/pwd")
	err = proc.Start()
	fmt.Println(err)
	proc.Stop()
	proc = NewExecutor(log, "donotexist")
	err = proc.Start()
	fmt.Println(err)
	proc.Stop()
	proc = NewExecutor(log, "/etc/passwd")
	err = proc.Start()
	fmt.Println(err)
	proc.Stop()
	// Output:
	// command exited
	// command exited
	// command exited
	// command exited
}

func ExampleNewExecutor_bc() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, "_test/bc.sh")
	err := proc.Start()
	fmt.Println(err)
	//proc.log <- true
	proc.io <- "2+2"
	fmt.Println(<-proc.io)
	// and now, exit detection
	proc.io <- "quit"
	proc.log <- true
	select {
	case in := <-proc.io:
		fmt.Println(in)
	case <-proc.exit:
		fmt.Println("exit")
	}
	waitabit()
	proc.Stop()
	dump(log)
	// Output:
	// <nil>
	// 4
	// exit
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func ExampleNewExecutor_hello() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, "_test/hello.sh")
	err := proc.Start()
	fmt.Println(err)
	proc.io <- `{"name":"Mike"}`
	fmt.Println(<-proc.io)
	proc.log <- true
	waitabit()
	proc.Stop()
	waitabit()
	_, ok := <-proc.io
	fmt.Printf("io %v\n", ok)
	dump(log)
	// Unordered output:
	// <nil>
	// {"hello": "Mike"}
	// msg=hello Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// io false
}

func ExampleNewExecutor_term() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, "_test/hello.sh")
	err := proc.Start()
	fmt.Println(err)
	proc.io <- `{"name":"*"}`
	var exited bool
	select {
	case <-proc.io:
		exited = false
	case <-proc.exit:
		exited = true
	}
	proc.log <- true
	fmt.Printf("exit %v\n", exited)
	waitabit()
	proc.Stop()
	waitabit()
	_, ok := <-proc.io
	fmt.Printf("io %v\n", ok)
	dump(log)
	// Unordered output:
	// <nil>
	// exit true
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// Goodbye!
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// io false
}
