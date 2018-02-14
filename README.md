# openwhisk-runtime-go

This is a (work in progress) OpenWhisk runtime for  Golang,  with replacement of the executable as the main server, instead of just invoking a process for each request.

![demo](./demo.png)

# Background

## How Go actions are currently implemented

Currently, Go actions in OpenWhisk are implemented using the generic Docker support. 

There is a python based server, listening for `/init` and `/run` requests. The `/init` will collect an executable and place in the current folder, while the `/run` will invoke the executable with `popen`, feeding the input and returning the output as log, and the last line as the result as a serialized json.

The problem is that spawning a new process for each request sound like a revival of the CGI.  It is certainly not the most efficient implementation.  Basically everyone moved from CGI executing processes to servers listening for requests since many many years ago.

Just for comparison, AWS Lambda supports Go implementing a server, listening and serving requests. 

## Why the exec?

The problem here is Python and Node are dynamic scripting languages, while Go is a compiled language.

Node and Python runtimes are both  servers, they receive the code of the function, “eval" the code and then execute it for serving requests. 

Go, generating an executable, cannot afford to do that. We cannot “eval” precompiled code. But it is also inefficient to spawn a new process for each function invocation. 

The solution here is to exec only once, when the runtime receive the executable of the function, at the `/init` time. 

Then you should replace the main executable and  serve the `/run` requests directly in the replaced executable. Of course this means that the replaced executable should be able to serve the /init requests too. All of this should go in a library

# How the new support works

The new support for Go will look like the following:

```
package main

import (
	"encoding/json"
	"fmt"

	"github.com/sciabarracom/openwhisk-runtime-go/openwhisk"
)

func hello(event json.RawMessage) (json.RawMessage, error) {
	// input and output
	var input struct{ Name string }
	var output struct {
		Greetings string `json:"greetings"`
	}
	// read the input event
	json.Unmarshal(event, &input)
	if input.Name != "" {
		// handle the event
		output.Greetings = "Hello, " + input.Name
		fmt.Println(output.Greetings)
		return json.Marshal(output)
	}
	return nil, fmt.Errorf("no name specified")
}

func main() {
	openwhisk.Start(hello)
}
```

The magic of serving `/init` and `/run` will live inside the library.

The `Start` function will start a web server listening for  the two requests of the proxy.

Posts to `/run` will invoke some json decoding  and then invoke the function.

Posts to `/init` will receive an executable, place somewhere `action` and then execute to it (expecting of course the server itself is implemented using the same library).  

# Testing the current implementation

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
OK
$ curl -XPOST http://localhost:8080/run -d '{"value":{"name":"Mike"}}'
{"greetings":"Hello, Mike"}
```

As you can see, the function changed and now it implements the "hello" handler.

But the replaced server is still able to run init so let's do it again, replacing with the "ciao" handler.


```
$ curl -XPOST http://localhost:8080/init -d @ciao.json
OK
$ curl -XPOST http://localhost:8080/run -d '{"value":{"name":"Mike"}}'
{"saluti":"Ciao, Mike"}
```



