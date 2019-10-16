package openwhisk

import (
	"io/ioutil"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, req *http.Request) {

	resp, err := http.Get("/hello")

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))

}
