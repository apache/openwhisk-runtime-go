#!/bin/bash
cd $(dirname $0)
bin/build.sh src/hello.go 
bin/build.sh src/hello_greeting.go 
bin/build.sh src/hello_message.go 
bin/build.sh src/empty.go
bin/build.sh src/hi.go
rm -Rvf action
#go run ../main/proxy.go #-debug
go build -o ../docker/proxy ../main/proxy.go
docker build -t golang19action ../docker/
docker run -ti -p 8080:8080  golang19action "$@" 
