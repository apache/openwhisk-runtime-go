package openwhisk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Params are the parameteres sent to the action
type Params struct {
	Value json.RawMessage `json:"value"`
}

func activationMessage() {
	fmt.Println("XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX")
	fmt.Println("XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX")
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

	// execute our action
	response, err := theAction(params.Value)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	// encode the response
	buf, err := json.Marshal(response)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error marshaling response: %v", err))
		return
	}

	// return the response
	w.Header().Set("Content-Type", "application/json")
	numBytesWritten, err := w.Write(buf)
	if err != nil {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Error writing response: %v", err))
		return
	}

	// diagnostic when writing problems
	if numBytesWritten != len(buf) {
		sendError(w, http.StatusInternalServerError, fmt.Sprintf("Only wrote %d of %d bytes to response", numBytesWritten, len(buf)))
		return
	}

}
