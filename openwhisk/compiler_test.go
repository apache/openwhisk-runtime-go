package openwhisk

import (
	"fmt"
	"io/ioutil"
)

/* this test confuses gogradle
func Example_compileAction_wrong() {
	sys("_test/precompile.sh", "hello.sh", "0")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/0", "../core/gobuild", log, true)
	fmt.Println(ap.CompileAction("_test/compile/0/exec", "exec"))
	// Output:
	// exit status 1
}*/

func Example_compileAction_singlefile_main() {
	sys("_test/precompile.sh", "hello.src", "1")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/1", "../core/gobuild", log, true)
	fmt.Println(ap.CompileAction("_test/compile/1/exec", "main"))
	sys("_test/postcompile.sh", "_test/compile/1/exec")
	// Output:
	// <nil>
	// _test/compile/1/exec: application/x-executable; charset=binary
	// name=Mike
	// {"message":"Hello, Mike!"}

}

func Example_compileAction_singlefile_hello() {
	sys("_test/precompile.sh", "hello1.src", "2")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/2", "../core/gobuild", log, true)
	fmt.Println(ap.CompileAction("_test/compile/2/exec", "hello"))
	sys("_test/postcompile.sh", "_test/compile/2/exec")
	// Output:
	// <nil>
	// _test/compile/2/exec: application/x-executable; charset=binary
	// name=Mike
	// {"hello":"Hello, Mike!"}
}

func Example_compileAction_multifile_main() {
	sys("_test/precompile.sh", "action", "3")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/3", "../core/gobuild", log, true)
	fmt.Println(ap.CompileAction("_test/compile/3/", "main"))
	sys("_test/postcompile.sh", "_test/compile/3/main")
	// Output:
	// <nil>
	// _test/compile/3/main: application/x-executable; charset=binary
	// Main:
	// Hello, Mike
	// {"greetings":"Hello, Mike"}
}

func Example_compileAction_multifile_hello() {
	sys("_test/precompile.sh", "action", "4")
	log, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy("./action/4", "../core/gobuild", log, true)
	fmt.Println(ap.CompileAction("_test/compile/4/", "hello"))
	sys("_test/postcompile.sh", "_test/compile/4/hello")
	// Output:
	// <nil>
	// _test/compile/4/hello: application/x-executable; charset=binary
	// Hello, Mike
	// {"greetings":"Hello, Mike"}

}
