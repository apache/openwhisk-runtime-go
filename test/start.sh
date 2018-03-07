#!/bin/bash
bin/build.sh ../main/hello.go 
bin/build.sh ../main/hello_greeting.go 
bin/build.sh ../main/hello_message.go 
bin/build.sh ../main/empty.go
bin/build.sh ../main/hi.go
rm -Rvf action
go run ../main/exec.go
