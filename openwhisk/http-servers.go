package openwhisk

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hello")

}
