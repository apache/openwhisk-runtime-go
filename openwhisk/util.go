package openwhisk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

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

func execActionIfExists() {
	action, err := exec.LookPath("./action")
	if err != nil {
		return
	}
	env := os.Environ()
	// shutdown the current server
	//theServer.Shutdown(nil)
	// execute the action
	err = syscall.Exec(action, nil, env)
	// restart the server if it failed
	theServer.ListenAndServe()
}
