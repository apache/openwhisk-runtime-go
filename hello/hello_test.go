package hello

import (
	"fmt"

	"github.com/sciabarracom/openwhisk-runtime-go/hello"
)

func ExampleHello() {
	name := []byte(`{ "name": "Mike"}`)
	data, _ := hello.Hello(name)
	fmt.Printf("%s", data)
	// Output:
	// {"greetings":"Hello, Mike"}
}

func ExampleHelloNoName() {
	name := []byte(`{ "noname": "Mike"}`)
	_, err := hello.Hello(name)
	fmt.Print(err)
	// Output:
	// no name specified
}
func ExampleHelloBadJson() {
	name := []byte(`{{`)
	_, err := hello.Hello(name)
	fmt.Print(err)
	// Output:
	// no name specified
}
