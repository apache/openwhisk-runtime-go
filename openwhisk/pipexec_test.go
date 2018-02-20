package openwhisk

import "fmt"

func ExampleNewPipeExec() {
	bc := NewPipeExec("bc", "-q")
	bc.print("2+2")
	fmt.Println(bc.scan())
	bc.print("3*3")
	fmt.Println(bc.scan())
	// Output:
	// 4
	// 9
}

func ExampleStartService() {
	ch := StartService("bc", "-q")
	ch <- "4+4"
	fmt.Println(<-ch)
	ch <- "8*8"
	fmt.Println(<-ch)
	// Output:
	// 8
	// 64
}
