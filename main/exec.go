package main

import (
	"encoding/json"
	"fmt"

	"github.com/sciabarracom/openwhisk-runtime-go/openwhisk"
)

func defaultAction(event json.RawMessage) (json.RawMessage, error) {
	return nil, fmt.Errorf("the action failed to locate a binary")
}

func main() {
	openwhisk.Start(defaultAction)
}
