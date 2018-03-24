package hello

import (
	"encoding/json"
	"fmt"
	"log"
)

// Hello receive an event in format
// { "name": "Mike"}
// and returns a greeting in format
// { "greetings": "Hello, Mike"}
func Hello(event json.RawMessage) (json.RawMessage, error) {
	// input and output
	var input struct {
		Name string
	}
	var output struct {
		Greetings string `json:"greetings"`
	}
	// read the input event
	json.Unmarshal(event, &input)
	if input.Name != "" {
		// handle the event
		output.Greetings = "Hello, " + input.Name
		log.Println(output.Greetings)
		return json.Marshal(output)
	}
	return nil, fmt.Errorf("no name specified")
}
