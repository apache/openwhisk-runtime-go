package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func hello(arg string) string {
	var obj map[string]interface{}
	json.Unmarshal([]byte(arg), &obj)
	name, ok := obj["name"].(string)
	if !ok {
		name = "Stranger"
	}
	msg := map[string]string{"message": ("Hello, " + name + "!")}
	res, _ := json.Marshal(msg)
	return string(res)
}

func main() {
	// native actions receive one argument, the JSON object as a string
	if len(os.Args) > 1 {
		fmt.Println(hello(os.Args[1]))
		return
	}
	// read loop
	fmt.Println(`{"openwhisk":1}`)
	reader := bufio.NewReader(os.Stdin)
	for {
		event, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Println(hello(event))
	}
}
