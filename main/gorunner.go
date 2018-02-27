package main

import (
	"fmt"

	"github.com/sciabarracom/openwhisk-runtime-go/openwhisk"
)

func hello() (string, error) {
	fmt.Println("Hello, world")
	return "Hello, world.", nil
	fm
}

func main() {
	openwhisk.Start(hello)
}
