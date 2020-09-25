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

var m = map[string]string{}

func ExampleNewExecutor_failed() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, log, "true", m)
	err := proc.Start(false)
	fmt.Println(err)
	proc.Stop()
	proc = NewExecutor(log, log, "/bin/pwd", m)
	err = proc.Start(false)
	fmt.Println(err)
	proc.Stop()
	proc = NewExecutor(log, log, "donotexist", m)
	err = proc.Start(false)
	fmt.Println(err)
	proc.Stop()
	proc = NewExecutor(log, log, "/etc/passwd", m)
	err = proc.Start(false)
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
	proc := NewExecutor(log, log, "_test/bc.sh", m)
	err := proc.Start(false)
	fmt.Println(err)
	res, _ := proc.Interact([]byte("2+2"))
	fmt.Printf("%s", res)
	proc.Stop()
	dump(log)
	// Output:
	// <nil>
	// 4
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func ExampleNewExecutor_hello() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, log, "_test/hello.sh", m)
	err := proc.Start(false)
	fmt.Println(err)
	res, _ := proc.Interact([]byte(`{"value":{"name":"Mike"}}`))
	fmt.Printf("%s", res)
	proc.Stop()
	dump(log)
	// Output:
	// <nil>
	// {"hello": "Mike"}
	// msg=hello Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func ExampleNewExecutor_env() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, log, "_test/env.sh", map[string]string{"TEST_HELLO": "WORLD", "TEST_HI": "ALL"})
	err := proc.Start(false)
	fmt.Println(err)
	res, _ := proc.Interact([]byte(`{"value":{"name":"Mike"}}`))
	fmt.Printf("%s", res)
	proc.Stop()
	dump(log)
	// Output:
	// <nil>
	// { "env": "TEST_HELLO=WORLD TEST_HI=ALL"}
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func ExampleNewExecutor_ack() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, log, "_test/hi", m)
	err := proc.Start(true)
	fmt.Println(err)
	proc.Stop()
	dump(log)
	// Output:
	// Command exited abruptly during initialization.
	// hi
}

func ExampleNewExecutor_badack() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, log, "_test/badack.sh", m)
	err := proc.Start(true)
	fmt.Println(err)
	proc.Stop()
	dump(log)
	// Output:
	// invalid character 'b' looking for beginning of value
}

func ExampleNewExecutor_badack2() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, log, "_test/badack2.sh", m)
	err := proc.Start(true)
	fmt.Println(err)
	proc.Stop()
	dump(log)
	// Output:
	// The action did not initialize properly.
}

func ExampleNewExecutor_helloack() {
	log, _ := ioutil.TempFile("", "log")
	proc := NewExecutor(log, log, "_test/helloack/exec", m)
	err := proc.Start(true)
	fmt.Println(err)
	res, _ := proc.Interact([]byte(`{"value":{"name":"Mike"}}`))
	fmt.Printf("%s", res)
	proc.Stop()
	dump(log)
	// Output:
	// <nil>
	// {"hello": "Mike"}
	// msg=hello Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}
