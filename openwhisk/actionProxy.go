package openwhisk

import (
	"log"
	"net/http"
)

// theServer is the current server
var theServer http.Server

// Start creates a proxy to execute actions
func Start() {

	// handle initialization
	http.HandleFunc("/init", initHandler)
	// handle execution
	http.HandleFunc("/run", runHandler)

	// start
	log.Println("Start!")
	theServer.Addr = ":8080"
	log.Fatal(theServer.ListenAndServe())
}
