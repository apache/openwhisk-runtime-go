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
	//log.Println("init: reading")
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
	}

	// decode request parameters
	//log.Println("init: decoding")
	var request initRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		sendError(w, http.StatusBadRequest, fmt.Sprintf("Error unmarshaling request: %v", err))
		return
	}

	// check if it is a binary
	if request.Value.Binary {
		var decoded []byte
		decoded, err = base64.StdEncoding.DecodeString(request.Value.Code)
		if err != nil {
			sendError(w, http.StatusBadRequest, "cannot decode the request: "+err.Error())
			return
		}
		// extract the replacement, stopping and then starting the action
		err = extractAction(&decoded, false)
	} else {
		buf := []byte(request.Value.Code)
		err = extractAction(&buf, true)
	}
	if err != nil {
		sendError(w, http.StatusBadRequest, "invalid action: "+err.Error())
		return
	}
	// stop and start
	reStartAction()
	if err != nil {
		sendError(w, http.StatusBadRequest, "cannot start action: "+err.Error())
		return
	}

	// answer OK
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", "3")
	w.Write([]byte("OK\n"))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

}
