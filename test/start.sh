#!/bin/bash
bin/build.sh src/hello.go 
bin/build.sh src/hello_greeting.go 
bin/build.sh src/hello_message.go 
bin/build.sh src/empty.go
bin/build.sh src/hi.go
rm -Rvf action
go run ../main/proxy.go #-debug
