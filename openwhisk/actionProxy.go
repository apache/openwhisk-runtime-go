package openwhisk

import (
	"log"
	"net/http"
	"os/exec"
)

// theServer is the current server
var theServer http.Server

// theChannel is the channel communicating with the action
var theChannel chan string

func stopAction() {
	// terminate current action
	if theChannel != nil {
		log.Println("terminating old action")
		theChannel <- ""
	}
}

func startAction() {
	// start a new action service
	_, err := exec.LookPath("./action/exec")
	if err == nil {
		log.Println("starting new action")
		theChannel = StartService("./action/exec")
	}
}

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
