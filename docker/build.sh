GOOS=linux GOARCH=amd64 go build -o proxy ../main/exec.go
cp proxy swift
docker build -t sciabarracom/openwhisk-runtime-go .
docker build -t sciabarracom/openwhisk-runtime-goswift swift
