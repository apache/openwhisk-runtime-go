go build -o hello ../main/hello.go
go build -o ciao ../main/ciao.go
echo '{"value":{"binary":true,"code":"'$(base64 hello)'"}}' >hello.json
echo '{"value":{"binary":true,"code":"'$(base64 ciao)'"}}' >ciao.json
go run ../main/exec.go
