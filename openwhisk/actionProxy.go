package openwhisk

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

// extractAction accept a byte array write it to a file
func extractAction(buf *[]byte) error {
	os.MkdirAll("./action", 0755)
	log.Println("Extract Action, assuming a binary")
	return ioutil.WriteFile("./action/exec", *buf, 0755)
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
