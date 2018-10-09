<!--
#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
-->

# Apache OpenWhisk Runtime for Go (and Generic executables)

[![Build Status](https://travis-ci.org/apache/incubator-openwhisk-runtime-go.svg?branch=master)](https://travis-ci.org/apache/incubator-openwhisk-runtime-go)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![Join Slack](https://img.shields.io/badge/join-slack-9B69A0.svg)](http://slack.openwhisk.org/)

:warning: Work in progress :warning:

# Apache OpenWhisk Runtime for Go (and Generic executables)

[![Build Status](https://travis-ci.org/apache/incubator-openwhisk-runtime-go.svg?branch=master)](https://travis-ci.org/apache/incubator-openwhisk-runtime-go)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![Join Slack](https://img.shields.io/badge/join-slack-9B69A0.svg)](http://slack.openwhisk.org/)

:warning: Work in progress :warning:

This is an OpenWhisk runtime for Golang and Generic executables.

- [Building it](#building)
- [Developing it](#development)
- [Using it with Go Sources](#gosources)
- [Precompiling Go Sources](#precompile)
- [Using it with Generic executables](#generic)

<a name="building"/>

# How to Build and Test

You need a Linux or an OSX environment, with Java and Docker installed to build the sources.

Prerequisites for running build and tests with gradle:

- docker
- jdk


To compile go proxy *in amd64 architecture* for docker:

```
./gradlew build
```

To build the docker images after compiling go proxy:

```
./gradlew distDocker
```

This will build the images:

* `actionloop-golang-v1.11`: an image supporting  Go sources
* `actionloop`: the base image, supporting generic executables ans shell script

The `actionloop` image can be used for supporting other compiled programming languages as long as they implement a `compile` script and the *action loop* protocol described below.

To run tests
```
./gradlew test --info
```
<a name="development"/>

# Local Development

If you want to develop the proxy and run tests natively, you can on Linux or OSX.
Tested on Ubuntu Linux (14.04) and OSX 10.13. Probably other distributions work, maybe even Windows with WSL, but since it is not tested YMMMV.

You need of course [go 1.11.0](https://golang.org/doc/install)

Then you need a set of utilities used in tests:

- bc
- zip
- realpath

Linux: `apt-get install bc zip realpath`
OSX: `brew install zip coreutils`

**NOTE**: Because tests build and cache some binary files, perform a `git clean -fx` and **do not share folders between linux and osx** because binaries are in different format...


<a name="gosources"/>

# Using it with Go Sources

The image can execute, compiling them on the fly, Golang OpenWhisk actions in source format.

An action must be a Go source file, placed in the `main` package and your action.

Since `main.main` is reserved in Golang for the entry point of your program, and the entry point is used by support code, your action must be named `Main` (with capital `M`) even if your specify `main` as the name of the action (or you do not specify it, defaulting to `main`). Also if you specify a function name different than `main`, for example `hello`, the name of your function  need to be capitalized, so your entry point will be `main.Hello`.

The expected signature for a `main` function is:

`func Main(event map[string]interface{}) map[string]interface{}`

For example:

```go
package main

import "log"

// Main is the function implementing the action
func Main(obj map[string]interface{}) map[string]interface{} {
  // do your work
  name, ok := obj["name"].(string)
  if !ok {
    name = "Stranger"
  }
  msg := make(map[string]interface{})
  msg["message"] = "Hello, " + name + "!"
  // log in stdout or in stderr
  log.Printf("name=%s\n", name)
  // encode the result back in json
  return msg
}
```

You can also have multiple source files in an action. In this case you need to collect all the sources  in a zip file for posting.

<a name="precompile"/>

## Precompiling Go Sources Offline

Compiling sources on the image can take some time when the images is initialised. You can speed up precompiling the sources using the image as an offline compiler. You need `docker` for doing that.

The images accepts a `-compile <main>` flag, and expects you provide sources in standard input. It will then compile them, emit the binary in standard output and errors in stderr. The output is always a zip file containing an executable.

If you have docker, you can do it this way:

If you have a single source maybe in file `main.go`, with a function named `Main` just do this:

`docker run openwhisk/actionloop-golang-v1.11 -compile main <main.go >main.zip`

If you have multiple sources in current directory, even with a subfolder with sources, you can compile it all with:

`zip -r - * | docker run openwhisk/actionloop-golang-v1.11 -compile main >main.zip`

The  generated executable is suitable to be deployed in OpenWhisk using just the generic `actionloop` runtime.

`wsk action create my/action main -docker openwhisk/actionloop`

You can also use the full `actionloop-golang-v1.11` as runtime, it is only bigger.

Note that the output is always a zip file in  Linux AMD64 format so the executable can be run only inside a Docker Linux container.

<a name="generic"/>

# Using it with generic Binaries

The `actionloop` image is designed to support generic Linux executable in an efficient way.

As such it works with any executable that supports the following simple protocol:

Repeat forever:
- read one line from standard input (file descriptor 0)
- parse the line as a json object, that will be in format:

```{
 "value": JSON,
 "namespace": String,
 "action_name": String,
 "api_host": String,
 "api_key": String,
 "activation_id": String,
 "deadline": Number
}```

Note that if you use libraries, those will expect the values in environment variables:

- `__OW_NAMESPACE`
- `__OW_ACTION_NAME`
- `__OW_API_HOST`
- `__OW_API_KEYS`
- `__OW_ACTIVATION_ID`
- `__OW_DEADLINE`

- execute the action, using the `value` that contains the payload provided by the used and logging in standard output and in standard error (file descriptor 1 and 2)
- write an answer in json format as a single line (without embedding newlines - newlines in strings must be quoted)

The `actionloop` image works actually with executable in unix sense, so also scripts are acceptable. In the actionloop image there is `bash` and the `jq` command, so you can for example implement the actionloop with a shell script:

```bash
#!/bin/bash
# read input forever line by line
while read line
do
   # parse the in input with `jq`
   name="$(echo $line | jq -r .name.value)"
   # log in stdout
   echo msg="hello $name"
   # produce the result - note the fd3
   echo '{"hello": "'$name'"}' >&3
done
```

Note the `actionloop` image will accept any source and will try to run it (if it is possible), while the `actionloop-golang`  images will try to compile the sources instead.


...Pending more detailed documentation...
