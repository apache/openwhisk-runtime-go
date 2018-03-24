package openwhisk

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartLatestAction(t *testing.T) {

	// cleanup
	os.RemoveAll("./action")
	theExecutor = nil

	// start an action that terminate immediately
	buf := []byte("#!/bin/sh\ntrue\n")
	extractAction(&buf, true)
	StartLatestAction()
	assert.Nil(t, theExecutor)

	// start the action that emits 1
	buf = []byte("#!/bin/sh\nwhile read a; do echo 1 >&3 ; done\n")
	extractAction(&buf, true)
	StartLatestAction()
	theExecutor.io <- "x"
	assert.Equal(t, <-theExecutor.io, "1")

	// now start an action that terminate immediately
	buf = []byte("#!/bin/sh\ntrue\n")
	extractAction(&buf, true)
	StartLatestAction()
	theExecutor.io <- "y"
	assert.Equal(t, <-theExecutor.io, "1")

	// start the action that emits 2
	buf = []byte("#!/bin/sh\nwhile read a; do echo 2 >&3 ; done\n")
	extractAction(&buf, true)
	StartLatestAction()
	theExecutor.io <- "z"
	assert.Equal(t, <-theExecutor.io, "2")
	/**/
	theExecutor.Stop()
}
