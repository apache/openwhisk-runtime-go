GOOS=linux GOARCH=amd64 go build -o proxy ../main/exec.go
docker build -t sciabarracom/openwhisk-exec .
GOOS=linux GOARCH=amd64 go build -o exec ../main/hello_exec.go
docker build -f Dockerfile.hello -t sciabarracom/openwhisk-hello .
