package openwhisk

import (
	"fmt"
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
	executable := fmt.Sprintf("./action/%d/exec", higherDir("./action"))
	_, err := exec.LookPath(executable)
	if err == nil {
		log.Printf("starting %s", executable)
		theChannel = StartService(executable)
	}
}

// Start creates a proxy to execute actions
func Start() {
	// handle initialization
	http.HandleFunc("/init", initHandler)
	// handle execution
	http.HandleFunc("/run", runHandler)
	// start action if there
	startAction()
	// start
	log.Println("Start!")
	theServer.Addr = ":8080"
	log.Fatal(theServer.ListenAndServe())
}
