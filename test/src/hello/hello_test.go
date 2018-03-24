package hello

import (
	"fmt"
)

func ExampleHello() {
	name := []byte(`{ "name": "Mike"}`)
	data, _ := Hello(name)
	fmt.Printf("%s", data)
	// Output:
	// {"greetings":"Hello, Mike"}
}

func ExampleHello_noName() {
	name := []byte(`{ "noname": "Mike"}`)
	_, err := Hello(name)
	fmt.Print(err)
	// Output:
	// no name specified
}
func ExampleHello_badJson() {
	name := []byte(`{{`)
	_, err := Hello(name)
	fmt.Print(err)
	// Output:
	// no name specified
}
