<!--
#
# Licensed to the Apache Software Foundation (ASF) under one or more contributor 
# license agreements.  See the NOTICE file distributed with this work for additional 
# information regarding copyright ownership.  The ASF licenses this file to you
# under the Apache License, Version 2.0 (the # "License"); you may not use this 
# file except in compliance with the License.  You may obtain a copy of the License 
# at:
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed 
# under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR 
# CONDITIONS OF ANY KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations under the License.
#
-->
# Apache OpenWhisk Runtime for Go (and Generic executables)

:warning: Work in progress :warning:

This is an OpenWhisk runtime for Golang and Generic executables.

- [Building it](#building)
- [Using it with Go Sources](#gosources)
- [Precompiling Go Sources](#precompile)
- [Using it with Generic executables](#generic)

<a name="building"/>

# How to Build

You need a linux environment, with Java and Docker installed to build the sources.

Prerequisites for running build and tests:
- docker
- jdk
- go 1.10.2
- bc (sudo apt-get install bc)

To compile go proxy
```
./gradlew build
```
To build the docker images after compiling go proxy
```
./gradlew distDocker
```

This will build the images:

* `actionloop-golang-v1.10`: an image supporting  Go sources
* `actionloop`: the base image, supporting generic executables

The `actionloop` image is used as a basis also for the `actionloop-swift` image. It can be used for supporting other compiled programming languages as long as they implement a `compile` script and the *action loop* protocol described below.

To run tests
```
./gradlew test --info
```

<a name="gosources"/>

# Using it with Go Sources

The image can execute, compiling them on the fly, Golang OpenWhisk actions in source format. An action must be a Go source file, placed in the `action` package, implementing the `Main` function (or the function specified as `main`).  

The expected signature is:

`func Main(event json.RawMessage) (json.RawMessage, error)`

Note the name of the function must be capitalised, because it needs to be exported from the `action` package. You can say the name of the function also in lower case, it will be capitalised anyway.

For example:

```
package action

import (
  "encoding/json"
  "log"
)

// Main is the function implementing the action
func Main(event json.RawMessage) (json.RawMessage, error) {
  // decode the json
  var obj map[string]interface{}
  json.Unmarshal(event, &obj)
  // do your work
  name, ok := obj["name"].(string)
  if !ok {
    name = "Stranger"
  }
  msg := map[string]string{"message": ("Hello, " + name + "!")}
  // log in stdout or in stderr 
  log.Printf("name=%s\n", name)
  // encode the result back in json
  return json.Marshal(msg)
}
```

You can also have multiple source files in an action. In this case you need to collect all the sources  in a zip file for posting.

<a name="precompile"/>

## Precompiling Go Sources Offline

Compiling sources on the image can take some time when the images is initialised. You can speed up precompiling the sources using the image as an offline compiler. You need `docker` for doing that.

The images accepts a `compile` command expecting sources in `/src`. It will then compile them and place the resut in `/out`.

If you have docker, you can do it this way:

- place your sources under `src` folder in current directory
- create an `out` folder to receive the binary
- run: `docker run -v $PWD/src:/src -v $PWD/out openwhisk/actionloop-golang-v1.10 compile`
- you can now use `wsk` to publish the `out/main` executable

If you have a function named in a different way, for example `Hello`, specify `compile hello`. It will produce a binary named `out/hello`

<a name="generic"/>

# Using it with generic Binaries

The `actionloop` image is designed to support generic linux executable in an efficient way. 

As such it works with any executable that supports the following simple protocol:

Repeat forever:
- read one line from stadard input (file descriptor 0)
- parse the line as a json object
- execute the action, logging in standard output and in standardar error (file descriptor 1 and 2)
- write an anwser in json format as a single line (without embedding newlines - newlines in strings must be quoted)

The `actionloop` image works actually with executable in unix sense, so also scripts are acceptable. In the actionloop image there is `bash` and the `jq` command, so you can for example implement the actionloop with a shell script:

```
#!/bin/bash
# read input forever line by line
while read line
do
   # parse the in input with `jq`
   name="$(echo $line | jq -r .name)"
   # log in stdout
   echo msg="hello $name"
   # produce the result - note the fd3
   echo '{"hello": "'$name'"}' >&3
done
```

Note the `actionloop` image will accept any source and will try to run it (if it is possible), while the `actionloop-golang` and `actionloop-swift` images will try to compile the sources instead.
