# openwhisk-runtime-go

This is a (work in progress) runtime for OpenWhisk for GO with replacement of the executable instead of executing them.

# Preparation

First, let's prepare the replacements:

```
cd test
go build -o hello ../main/hello.go
go build -o ciao ../main/ciao.go
echo '{"value":{"binary":true,"code":"'$(base64 hello)'"}}' >hello.json
echo '{"value":{"binary":true,"code":"'$(base64 ciao)'"}}' >ciao.json
```

Now, start the server:

```
go run ../main/exec.go
```

# You can now test the hello functions

Default behaviour (no executable)

```
$ curl -XPOST http://localhost:8080/run -d '{"value":{"name":"Mike"}}'
{"error":"the action failed to locate a binary"}
```

Now post the `hello` handler and run it:

```
$ curl -XPOST http://localhost:8080/init -d @hello.json
$ curl -XPOST http://localhost:8080/run -d '{"value":{"name":"Mike"}}'
```

Now post the `ciao` handler and run it:

```
$ curl -XPOST http://localhost:8080/init -d @ciao.json
$ curl -XPOST http://localhost:8080/run -d '{"value":{"name":"Mike"}}'
```



