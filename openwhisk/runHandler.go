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

// ErrResponse is the response when there are errors
type ErrResponse struct {
	Error string `json:"error"`
}

func sendError(w http.ResponseWriter, code int, cause string) {
	fmt.Println("action error:", cause)
	errResponse := ErrResponse{Error: cause}
	b, err := json.Marshal(errResponse)
	if err != nil {
		fmt.Println("error marshalling error response:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(b)
}

func runHandler(w http.ResponseWriter, r *http.Request) {

	// reading the request
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
	response, err := Action(params.Value)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error: %v", err))
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
