package openwhisk

import (
	"reflect"
)

// Start creates a proxy to execute actions
func Start(actionSymbol interface{}) {
	action := reflect.ValueOf(actionSymbol)
	var args []reflect.Value
	action.Call(args)
}
