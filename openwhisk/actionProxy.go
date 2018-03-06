package openwhisk

import (
	"fmt"
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

func reStartAction() {
	// stop action if any
	stopAction()
	// find the action if any
	highestDir := highestDir("./action")
	if highestDir == 0 {
		log.Println("no action dir")
		theChannel = nil
		return
	}
	// try to launch the action
	executable := fmt.Sprintf("./action/%d/exec", highestDir)
	_, err := exec.LookPath(executable)
	// try to start the action
	if err == nil {
		log.Printf("starting %s", executable)
		ch := StartService(executable)
		if ch != nil {
			theChannel = ch
			return
		}
	}
	// cannot start, removing the action and retry
	exeDir := fmt.Sprintf("./action/%d/", highestDir)
	os.RemoveAll(exeDir)
	reStartAction()
}

// Start creates a proxy to execute actions
func Start() {
	// handle initialization
	http.HandleFunc("/init", initHandler)
	// handle execution
	http.HandleFunc("/run", runHandler)
	// start action if there
	reStartAction()
	// start
	log.Println("Started!")
	theServer.Addr = ":8080"
	log.Fatal(theServer.ListenAndServe())
}
