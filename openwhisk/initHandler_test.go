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

import "path/filepath"

func Example_badinit_nocompiler() {
	ts, cur, log := startTestServer("")
	sys("_test/build.sh")
	doRun(ts, "")
	doInit(ts, "{}")
	//sys("ls", "_test/exec")
	doInit(ts, initBinary("_test/exec", ""))      // empty
	doInit(ts, initBinary("_test/hi", ""))        // say hi
	doInit(ts, initBinary("_test/hello.src", "")) // source not excutable
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 400 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 400 {"error":"cannot start action: command exited"}
	// 400 {"error":"cannot start action: command exited"}
	// 400 {"error":"cannot start action: command exited"}
	// 400 {"error":"no action defined yet"}

}

func Example_bininit_nocompiler() {
	ts, cur, log := startTestServer("")
	sys("_test/build.sh")
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_message", ""))
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_greeting", ""))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 400 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"message":"Hello, Mike!"}
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_zipinit_nocompiler() {
	ts, cur, log := startTestServer("")
	sys("_test/build.sh")
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_greeting.zip", ""))
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_message.zip", ""))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 400 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// 200 {"ok":true}
	// 200 {"message":"Hello, Mike!"}
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

/* commented out test for timing problems
func Example_shell_nocompiler() {
	ts, cur, log := startTestServer("")
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello.sh", ""))
	doRun(ts, "")
	doRun(ts, `{"name":"*"}`)
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 400 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"hello": "Mike"}
	// 400 {"error":"command exited"}
	// 400 {"error":"no action defined yet"}
	// msg=hello Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// Goodbye!
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}*/

func Example_main_nocompiler() {
	ts, cur, log := startTestServer("")
	sys("_test/build.sh")
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_message", "message"))
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_greeting", "greeting"))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 400 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"message":"Hello, Mike!"}
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_main_zipinit_nocompiler() {
	ts, cur, log := startTestServer("")
	sys("_test/build.sh")
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_greeting.zip", "greeting"))
	doInit(ts, initBinary("_test/hello_greeting1.zip", "greeting"))
	doRun(ts, "")
	doInit(ts, initBinary("_test/hello_message.zip", "message"))
	doInit(ts, initBinary("_test/hello_message1.zip", "message"))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 400 {"error":"no action defined yet"}
	// 400 {"error":"cannot start action: command exited"}
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// 400 {"error":"cannot start action: command exited"}
	// 200 {"ok":true}
	// 200 {"message":"Hello, Mike!"}
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_compile_simple() {
	comp, _ := filepath.Abs("../core/gobuild")
	ts, cur, log := startTestServer(comp)
	sys("_test/build.sh")
	doRun(ts, "")
	doInit(ts, initCode("_test/hello.src", ""))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 400 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"message":"Hello, Mike!"}
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_compile_withMain() {
	comp, _ := filepath.Abs("../core/gobuild")
	ts, cur, log := startTestServer(comp)
	sys("_test/build.sh")
	doRun(ts, "")
	doInit(ts, initCode("_test/hello1.src", ""))
	doInit(ts, initCode("_test/hello1.src", "hello"))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 400 {"error":"no action defined yet"}
	// 400 {"error":"cannot start action: command exited"}
	// 200 {"ok":true}
	// 200 {"hello":"Hello, Mike!"}
	// name=Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

func Example_compile_withZipSrc() {
	sys("_test/zips.sh")
	comp, _ := filepath.Abs("../core/gobuild")
	ts, cur, log := startTestServer(comp)
	doRun(ts, "")
	doInit(ts, initBinary("_test/action.zip", ""))
	doRun(ts, "")
	doInit(ts, initBinary("_test/action.zip", "hello"))
	doRun(ts, "")
	stopTestServer(ts, cur, log)
	// Output:
	// 400 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// 200 {"ok":true}
	// 200 {"greetings":"Hello, Mike"}
	// Main:
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// Hello, Mike
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}

/*
func Example_compile_withZipSrcDefault() {
	sys("_test/zips.sh")
	comp, _ := filepath.Abs("../core/gobuild")
	ts, cur := startTestServer(comp)
	doRun(ts, "")
	doInit(ts, initBinary("_test/action.zip", ""))
	doRun(ts, "")
	stopTestServer(ts, cur)
	// Output:
	// 400 {"error":"no action defined yet"}
	// 200 {"ok":true}
	// name=Mike
	// 200 {"hello":"Hello, Mike!"}
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
	// XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX
}
/**/
