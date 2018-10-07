package wskcli

import (
	"net/http"

	"github.com/apache/incubator-openwhisk-client-go/whisk"
)

// Invoke a remote action
func Invoke(actionName string, payload interface{}) (map[string]interface{}, error) {
	client, err := whisk.NewClient(http.DefaultClient, nil)
	if err != nil {
		return nil, err
	}
	resp, _, err := client.Actions.Invoke(actionName, payload, true, true)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
