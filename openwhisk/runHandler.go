package openwhisk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

// theChannel is the channel communicating with the action
var theChannel chan string

// Params are the parameteres sent to the action
type Params struct {
	Value json.RawMessage `json:"value"`
}

func activationMessage() {
	fmt.Println("XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX")
	fmt.Println("XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX")
}

// ErrResponse is the response when there are errors
type ErrResponse struct {
	Error string `json:"error"`
}

func sendError(w http.ResponseWriter, code int, cause string) {
	fmt.Println(cause)
	errResponse := ErrResponse{Error: cause}
	b, err := json.Marshal(errResponse)
	if err != nil {
		fmt.Println("error marshalling error response:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func startActionIfExists() {
	// terminate current action
	if theChannel != nil {
		log.Println("terminating old action")
		theChannel <- ""
	}
	// start a new action service
	_, err := exec.LookPath("./action")
	if err == nil {
		log.Println("starting new action")
		theChannel = StartService("./action")
	}
}

func runHandler(w http.ResponseWriter, r *http.Request) {

	// send the activation message at the end
	defer activationMessage()

	// parse the request
	params := Params{}
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error reading request body: %v", err))
		return
	}

	// decode request parameters
	err = json.Unmarshal(body, &params)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error unmarshaling request: %v", err))
		return
	}

	// check if you have an action
	if theChannel == nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("no action defined yet"))
		return
	}

	// execute the action
	theChannel <- string(params.Value)
	response := <-theChannel
	if response == "" {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	// return the response
	w.Header().Set("Content-Type", "application/json")
	numBytesWritten, err := w.Write([]byte(response))
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error writing response: %v", err))
		return
	}

	// diagnostic when writing problems
	if numBytesWritten != len(response) {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Only wrote %d of %d bytes to response", numBytesWritten, len(response)))
		return
	}
}
