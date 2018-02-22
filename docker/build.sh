GOOS=linux GOARCH=amd64 go build -o proxy ../main/exec.go
docker build -f Dockerfile.go -t sciabarracom/openwhisk-runtime-go .
docker build -f Dockerfile.goswift -t sciabarracom/openwhisk-runtime-goswift .
