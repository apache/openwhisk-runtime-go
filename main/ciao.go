package main

import (
	"encoding/json"
	"fmt"

	"github.com/sciabarracom/openwhisk-runtime-go/openwhisk"
)

func ciao(event json.RawMessage) (json.RawMessage, error) {
	// input and output
	var input struct{ Name string }
	var output struct {
		Greetings string `json:"saluti"`
	}
	// read the input event
	json.Unmarshal(event, &input)
	if input.Name != "" {
		// handle the event
		output.Greetings = "Ciao, " + input.Name
		fmt.Println(output.Greetings)
		return json.Marshal(output)
	}
	return nil, fmt.Errorf("no name specified")
}

func main() {
	openwhisk.Start(ciao)
}
