package openwhisk

import (
	"io/ioutil"
	"log"
	"os"
)

// extractAction accept a byte array write it to a file
func extractAction(buf *[]byte) error {
	os.MkdirAll("./action", 0755)
	log.Println("Extract Action, assuming a binary")
	return ioutil.WriteFile("./action/exec", *buf, 0755)
}
