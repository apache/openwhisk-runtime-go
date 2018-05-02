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

/* this test confuses gogradle
func Example_compileAction_wrong() {
	sys("_test/precompile.sh", "hello.sh", "0")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/0", "../common/gobuild.sh", log)
	fmt.Println(ap.CompileAction("_test/compile/0/exec", "exec", ""))
	// Output:
	// exit status 1
}*/

func Example_isCompiled() {
	sys("_test/precompile.sh", "hello.src", "c")
	file := abs("./_test/compile/c/exec")
	dir := abs("./_test/compile/c")
	fmt.Println(isCompiled(file, "main"))
	fmt.Println(isCompiled(dir, "exec"))

	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/c", "../common/gobuild.sh", log)
	ap.CompileAction("main", abs("_test/compile/c/exec"), "")

	fmt.Println(isCompiled(file, "main"))
	fmt.Println(isCompiled(dir, "exec"))
	// errors
	fmt.Println(isCompiled(dir, "main"))
	fmt.Println(isCompiled(file+"1", "main"))

	// Output:
	// false
	// false
	// true
	// true
	// false
	// false
}

func Example_compileAction_singlefile_main() {
	sys("_test/precompile.sh", "hello.src", "1")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/1", "../common/gobuild.sh", log)
	fmt.Println(ap.CompileAction("main", abs("_test/compile/1/exec"), ""))
	sys("_test/postcompile.sh", "_test/compile/1/exec")
	// Output:
	// <nil>
	// _test/compile/1/exec: application/x-executable; charset=binary
	// name=Mike
	// {"message":"Hello, Mike!"}
}

func Example_compileAction_singlefile_main_out() {
	sys("_test/precompile.sh", "hello.src", "1a")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/1a", "../common/gobuild.sh", log)
	fmt.Println(ap.CompileAction("main", abs("_test/compile/1a/exec"), abs("_test/output/1a")))
	sys("_test/postcompile.sh", "_test/output/1a/main")
	// Output:
	// <nil>
	// _test/output/1a/main: application/x-executable; charset=binary
	// name=Mike
	// {"message":"Hello, Mike!"}
}

func Example_compileAction_singlefile_hello() {
	sys("_test/precompile.sh", "hello1.src", "2")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/2", "../common/gobuild.sh", log)
	fmt.Println(ap.CompileAction("hello", "_test/compile/2/exec", ""))
	sys("_test/postcompile.sh", "_test/compile/2/exec")
	// Output:
	// <nil>
	// _test/compile/2/exec: application/x-executable; charset=binary
	// name=Mike
	// {"hello":"Hello, Mike!"}
}

func Example_compileAction_singlefile_hello_out() {
	sys("_test/precompile.sh", "hello1.src", "2a")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/2a", "../common/gobuild.sh", log)
	fmt.Println(ap.CompileAction("hello", "_test/compile/2a/exec", abs("_test/output/2a")))
	sys("_test/postcompile.sh", "_test/output/2a/hello")
	// Output:
	// <nil>
	// _test/output/2a/hello: application/x-executable; charset=binary
	// name=Mike
	// {"hello":"Hello, Mike!"}
}

func Example_compileAction_multifile_main() {
	sys("_test/precompile.sh", "action", "3")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/3", "../common/gobuild.sh", log)
	fmt.Println(ap.CompileAction("main", "_test/compile/3/", ""))
	sys("_test/postcompile.sh", "_test/compile/3/main")
	// Output:
	// <nil>
	// _test/compile/3/main: application/x-executable; charset=binary
	// Main:
	// Hello, Mike
	// {"greetings":"Hello, Mike"}
}

func Example_compileAction_multifile_main_out() {
	sys("_test/precompile.sh", "action", "3a")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/3a", "../common/gobuild.sh", log)
	fmt.Println(ap.CompileAction("main", "_test/compile/3a/", abs("_test/output/3a")))
	sys("_test/postcompile.sh", "_test/output/3a/main")
	// Output:
	// <nil>
	// _test/output/3a/main: application/x-executable; charset=binary
	// Main:
	// Hello, Mike
	// {"greetings":"Hello, Mike"}
}

func Example_compileAction_multifile_hello() {
	sys("_test/precompile.sh", "action", "4")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/4", "../common/gobuild.sh", log)
	fmt.Println(ap.CompileAction("hello", "_test/compile/4/", ""))
	sys("_test/postcompile.sh", "_test/compile/4/hello")
	// Output:
	// <nil>
	// _test/compile/4/hello: application/x-executable; charset=binary
	// Hello, Mike
	// {"greetings":"Hello, Mike"}
}

func Example_compileAction_multifile_hello_out() {
	sys("_test/precompile.sh", "action", "4a")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/4a", "../common/gobuild.sh", log)
	fmt.Println(ap.CompileAction("hello", "_test/compile/4/", abs("_test/output/4a")))
	sys("_test/postcompile.sh", "_test/output/4a/hello")
	// Output:
	// <nil>
	// _test/output/4a/hello: application/x-executable; charset=binary
	// Hello, Mike
	// {"greetings":"Hello, Mike"}
}
