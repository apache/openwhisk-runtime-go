#!/bin/bash
bin/build.sh ../main/hello_greeting.go 
bin/build.sh ../main/hello_message.go 
rm -Rvf action
go run ../main/exec.go
