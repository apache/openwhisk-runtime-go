package openwhisk

import "net/http"

func (ap *ActionProxy) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello\n"))
}
