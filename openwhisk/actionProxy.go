package openwhisk

import (
	"encoding/json"
	"log"
	"net/http"
)

// Action is the actual action executed
var Action func(json.RawMessage) (json.RawMessage, error)

// Start creates a proxy to execute actions
func Start(action func(json.RawMessage) (json.RawMessage, error)) {
	// set the action at the start
	Action = action
	// handle initialization
	http.HandleFunc("/init", initHandler)
	// handle execution
	http.HandleFunc("/run", runHandler)
	// start
	log.Fatal(http.ListenAndServe(":8080", nil))
}
