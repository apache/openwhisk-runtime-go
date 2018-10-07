package main

import (
	"wskcli"
)

func mkerr(err error) map[string]interface{} {
	res := make(map[string]interface{})
	res["error"] = err.Error()
	return res
}

// Hello is the function implementing the action
func Hello(obj map[string]interface{}) map[string]interface{} {
	// encode the result back in json
	resp, err := wskcli.Invoke(obj["action"].(string), obj)
	if err != nil {
		return mkerr(err)
	}
	return resp
}
