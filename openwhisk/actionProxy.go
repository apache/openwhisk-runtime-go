package openwhisk

import (
	"encoding/json"
	"log"
	"net/http"
)

// theAction is the actual action executed
var theAction func(json.RawMessage) (json.RawMessage, error)

// theServer is the current server
var theServer http.Server

// Start creates a proxy to execute actions
func Start(action func(json.RawMessage) (json.RawMessage, error)) {

	// set the action at the start
	theAction = action
	// handle initialization
	http.HandleFunc("/init", initHandler)
	// handle execution
	http.HandleFunc("/run", runHandler)

	// start
	theServer.Addr = ":8080"
	log.Fatal(theServer.ListenAndServe())
}
