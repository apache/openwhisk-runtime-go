package action

import (
	"encoding/json"
	"fmt"
)

// Main forwading to Hello
func Main(event json.RawMessage) (json.RawMessage, error) {
	fmt.Println("Main:")
	return Hello(event)
}
