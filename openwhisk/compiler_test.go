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
	"os"
)

/**
Notes to understand tests:
- tests are run from the "openwhisk" folder, as the current directory
- precompile.sh prepare a compilation environment:
	_test/precompile.sh hello.src aaa main
  produces
	 - _test/compile/src/aaa/src/main
	 - _test/compile/src/aaa/bin/
  ready for the compiler
- compiler (../../common/gobuild.py) takes 3 arguments:
   <main> <source-dir> <target-dir>
   - it will look for a <source-dir>/<main> file
   - will generate some files in <source-dir>
   - compiler output is in <target-dir>/<main>
 - postcompile.sh will
	- execute the binary with 3>&1
	- feed it with the json '{"name":"Mike"}
	- will print the type of the executable and its output and log
*/

const (
	PREP  = "_test/precompile.sh"
	CHECK = "_test/postcompile.sh"
	TMP   = "_test/compile/"
	COMP  = "../common/gobuild.py"
)

// compile a main
func Example_cli_compiler() {
	sys(PREP, "hello.src", "0", "exec")
	ap := NewActionProxy(TMP, COMP, os.Stdout, os.Stderr)
	fmt.Println(isCompiled(TMP + "0/src/exec"))
	ap.CompileAction("main", TMP+"0/src", TMP+"0/bin")
	sys(CHECK, TMP+"0/bin/exec")
	fmt.Println(isCompiled(TMP + "0/bin/exec"))
	// errors
	fmt.Println(isCompiled(TMP + "0/bin/exec1"))
	// Output:
	// false
	// _test/compile/0/bin/exec: application/x-executable
	// name=Mike
	// {"message":"Hello, Mike!"}
	// true
	// false
}

// compile a not-main (hello) function
func Example_hello() {
	N := "1"
	sys(PREP, "hello1.src", N, "exec")
	ap := NewActionProxy(TMP, COMP, os.Stdout, os.Stderr)
	env := map[string]interface{}{"GOROOT": TMP + N}
	ap.SetEnv(env)
	ap.CompileAction("hello", TMP+N+"/src", TMP+N+"/bin")
	sys(CHECK, TMP+N+"/bin/exec")
	// Output:
	// _test/compile/1/bin/exec: application/x-executable
	// name=Mike
	// {"hello":"Hello, Mike!"}
}

// compile a function including a package
func Example_package() {
	N := "2"
	sys(PREP, "hello2.src", N, "exec", "hello")
	ap := NewActionProxy(TMP, COMP, os.Stdout, os.Stderr)
	env := map[string]interface{}{"GOROOT": TMP + N}
	ap.SetEnv(env)
	ap.CompileAction("main", TMP+N+"/src", TMP+N+"/bin")
	sys(CHECK, TMP+N+"/bin/exec")
	// Output:
	// _test/compile/2/bin/exec: application/x-executable
	// Main
	// Hello, Mike
	// {"greetings":"Hello, Mike"}
}

func Example_compileError() {
	N := "6"
	sys(PREP, "error.src", N)
	ap := NewActionProxy(TMP, COMP, os.Stdout, os.Stderr)
	err := ap.CompileAction("main", TMP+N+"/src", TMP+N+"/bin")
	fmt.Printf("%v", removeLineNr(err.Error()))
	// Unordered output:
	// ./exec__.go::: syntax error: unexpected error at end of statement
}

func Example_withMain() {
	N := "7"
	sys(PREP, "hi.src", N, "exec")
	ap := NewActionProxy(TMP, COMP, os.Stdout, os.Stderr)
	err := ap.CompileAction("main", TMP+N+"/src", TMP+N+"/bin")
	fmt.Println(err)
	sys(TMP + N + "/bin/exec")
	// Output:
	// <nil>
	// hi
}
