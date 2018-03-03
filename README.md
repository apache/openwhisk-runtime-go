OUTDATED README - currently actions are implemented with a pipe - will update soon.

# openwhisk-runtime-go
This is a (work in progress) OpenWhisk runtime for  Golang.

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


# Immplementation

TODO: describe here the new implementation with a piped child process reading input and output
