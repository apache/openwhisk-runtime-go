package hello

import (
	"encoding/json"
	"fmt"
)

// Ciao receive an event in format
// { "name": "Mike"}
// and returns a greeting in format
// { "saluti": "Ciao, Mike"}
func Ciao(event json.RawMessage) (json.RawMessage, error) {
	// input and output
	var input struct {
		Name string
	}
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
