package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sciabarracom/openwhisk-runtime-go/hello"
)

func main() {
	// handle command line argument
	if len(os.Args) > 1 {
		result, err := hello.Hello([]byte(os.Args[1]))
		if err == nil {
			fmt.Println(string(result))
			return
		}
		fmt.Printf("{ error: %q}\n", err.Error())
		return
	}
	// read loop
	fmt.Println(`{"openwhisk":1}`)
	reader := bufio.NewReader(os.Stdin)
	for {
		event, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		result, err := hello.Hello(event)
		if err != nil {
			fmt.Printf("{ error: %q}\n", err.Error())
			continue
		}
		fmt.Println(string(result))
	}
}
