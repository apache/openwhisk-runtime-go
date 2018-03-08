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

import "fmt"

func ExampleNewPipeExec() {
	bc := NewPipeExec("_test/bc.sh")
	bc.print("2+2")
	fmt.Println(bc.scan())
	bc.print("3*3")
	fmt.Println(bc.scan())
	// Output:
	// 4
	// 9
}

func ExampleNewPipeExec_failed() {
	proc := NewPipeExec("true")
	fmt.Println(proc.err)
	proc = NewPipeExec("pwd")
	fmt.Println(proc.err == nil)
	// Output:
	// no handshake
	// false
}

func ExampleStartService() {
	ch := StartService("_test/bc.sh")
	ch <- "4+4"
	fmt.Println(<-ch)
	ch <- "8*8"
	fmt.Println(<-ch)
	// Output:
	// 8
	// 64
}

func ExampleStartService_donotexistexit() {
	// do not exist
	ch := StartService("donotexist")
	fmt.Println(ch)
	// not a binary
	ch = StartService("/etc/passwd")
	fmt.Println(ch)
	ch = StartService("/bin/pwd")
	fmt.Println(ch)
	ch = StartService("true")
	fmt.Println(ch)
	// Output:
	// <nil>
	// <nil>
	// <nil>
	// <nil>
}

func ExampleStartService_exit() {
	ch := StartService("_test/bc.sh")
	if ch != nil {
		fmt.Println("channel not nil")
	}
	ch <- "4+4"
	_, ok := <-ch
	fmt.Println(ok)
	ch <- "quit"
	_, ok = <-ch
	fmt.Println(ok)
	// Output:
	// channel not nil
	// true
	// false
}

func ExampleStartService_true() {
	ch := StartService("/bin/pwd")
	fmt.Println(ch)
	// Output:
	// <nil>
}
