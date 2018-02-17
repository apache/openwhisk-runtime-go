package main

import (
	"github.com/sciabarracom/openwhisk-runtime-go/hello"
	"github.com/sciabarracom/openwhisk-runtime-go/openwhisk"
)

func main() {
	openwhisk.Start(hello.Hello)
}
