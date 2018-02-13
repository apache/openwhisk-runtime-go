package openwhisk

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type initRequest struct {
	Value struct {
		Code   string
		Binary bool `json:",omitempty"`
	}
}

func initHandler(w http.ResponseWriter, r *http.Request) {

	// read body of the request
	fmt.Println("init: reading")
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	// decode request parameters
	fmt.Println("init: decoding")
	var request initRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error unmarshaling request: %v", err))
		return
	}

	//request.Value.Code = fmt.Sprintf("%d", len(request.Value.Code))
	//fmt.Printf("init: req %v", request)

	// check if it is a binary
	if !request.Value.Binary {
		sendError(w, http.StatusBadRequest, "Source requests not (yet) supported")
		return
	}

	// write the base64 encoded file
	decoded, err := base64.StdEncoding.DecodeString(request.Value.Code)
	if err != nil {
		sendError(w, http.StatusBadRequest, "cannot decode the request")
		return
	}
	err = ioutil.WriteFile("./action", decoded, 0755)
	if err != nil {
		sendError(w, http.StatusBadRequest, "cannot write the file")
		return
	}
	fmt.Println("init: file written")

	// answer OK
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "3")
	w.Write([]byte("OK\n"))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// I've seen things you people wouldn't believe.
	// Attack ships on fire off the shoulder of Orion.
	// I watched C-beams glitter in the dark near the Tannhäuser Gate.
	// All those moments will be lost in time, like tears in rain. Time to die.
	execActionIfExists()
}
