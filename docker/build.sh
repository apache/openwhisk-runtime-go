GOOS=linux GOARCH=amd64 go build -o proxy ../main/exec.go
docker build -t sciabarracom/openwhisk-runtime-go .
