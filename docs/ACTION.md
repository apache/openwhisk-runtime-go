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
# Writing actions for the Golang and ActionLoop runtime

<a name="golang"/>

## How to write Go Actions

The `action-golang-v1.14` runtime can execute actions written in the Go programming language in OpenWhisk, either precompiled binary or compiling sources on the fly.

### Entry Point

The source of one action is one or more Go source file. The entry point of the action is a function, placed in the `main` package. The obvious name for the default action would be `main`, but unfortunately `main.main` is the fixed entry point, for a go program, and its signature is `main()` (without arguments) so it cannot be used, unless you implement the [ActionLoop](#actionloopgo) directly, overwriting the one provided by the runtime, see below.

When deploying an OpenWhisk action you can specify the `main` function, and the default value is of course `main`.

The rule used by the runtime you use the  the *capitalized* name of the function specified as main. The default is of course `main.Main`' if you specify `hello` it will be `hello.Hello`. It will be also `main.Main` or `hello.Hello` if you specify the main function as, respectively, `Main` and `Hello`. The function must have a specific signature, as described next.

*NOTE* The runtime does *not* support different packages from `main` for the entry point. If you specify `hello.main` the runtime will try to use `Hello.main`, that will be almost certainly incorrect. You can however have other packages in your sources, as described below.

### Signature

The expected signature for a `main` function is:

`func Main(event map[string]interface{}) map[string]interface{}`

So a very simple `hello world` function would be:

```go
package main

import "log"

// Main is the function implementing the action
func Main(obj map[string]interface{}) map[string]interface{} {
  // do your work
  name, ok := obj["name"].(string)
  if !ok {
    name = "world"
  }
  msg := make(map[string]interface{})
  msg["message"] = "Hello, " + name + "!"
  // log in stdout or in stderr
  log.Printf("name=%s\n", name)
  // encode the result back in json
  return msg
}
```

For the return result, not only support `map[string]interface{}` but also support `[]interface{}`

So a very simple `hello array` function would be:

```go
package main

// Main is the function implementing the action
func Main(event map[string]interface{}) []interface{} {
        result := []interface{}{"a", "b"}
        return result
}
```

And support array result for sequence action as well, the first action's array result can be used as next action's input parameter.

So the function can be:

```go
package main

// Main is the function implementing the action
func Main(obj []interface{}) []interface{} {
        return obj
}
```

You can also have multiple source files in an action, packages and vendor folders.  Check the [deployment](DEPLOY.md) document for more details how to package and deploy actions.

<a name="generic"/>

## Using it with generic Binaries

The `actionloop` runtime can execute  generic Linux executable in an efficient way. The actions should work reading input line by line, perform its work and produce output also line by line. In more detail it should respect the following protocol.

<a name="actionloop">

### The Action Loop Protocol

The protocol can be specified informally as follows.

- Send an acknowledgement after initialization when required. If the environment variable `__OW_WAIT_FOR_ACK` is not empty, write on file descriptor 3 the string `{ "ok": true }`.
- Read one line from standard input (file descriptor 0).
- Parse the line as a JSON object. Currently the object will be in currently in the format:

```
{
 "value": JSON,
 "namespace": String,
 "action_name": String,
 "api_host": String,
 "api_key": String,
 "activation_id": String,
 "deadline": Number
}
```

Note however that more values could be provided in future.
Usually this JSON is read and the values are stored in environment variables, converted to upper case the key and  and adding the prefix `__OW_`.

- The payload of the request is stored in the key `value`. The action should read the field `value` assuming it is a JSON object (note, not an array, nor a string or number) and parse it.
- The action can now perform its tasks as appropriate. The action can produce log writing  in standard output (file descriptor 1) and standard error (file descriptor 3) . Note that those corresponds to file descriptors 1 and 2.
- The action will receive also file descriptor 3 for returning results. The result of the action must be a single line (without embedding newlines - newlines in strings must be quoted) written in file descriptor 3.
- The action should not exit now, but continue the loop, reading the next line and processing as described before, continuing forever.

### Using shell scripts

The `actionloop` image works actually with executable in Linux sense, so also scripts are acceptable.

In the current actionloop image there is `bash` and the `jq` command, so you can for example implement the actionloop with a shell script like this:

```bash
#!/bin/bash
# send an ack if required
if test -n "$__OW_WAIT_FOR_ACK"
  then echo '{"ok":true}' >&3
fi
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

Note here we are just interested in the payload, but in general you may also want to retrieve other fields.

Note the `actionloop` image will accept any source and will try to run it (if it is possible), while the `action-golang-v1.N`  will instead try to compile the sources assuming it is Golang instead.

<a name="actionloopgo">

### Providing your own ActionLoop implementation

By default the runtime expects you provide a main function that will serve one request, and will add a default implementation of the ActionLoop protocol when compiling.

You can however overwrite the default protocol and provide your how implementation of the ActionLoop. If you do so, you will have to take care of opening file descriptors, reading input, parse JSON and set environment variables.

To overwrite the default ActionLoop you can do this either sending a single file source action, or a zip action.

If you send a single file, you have to provide your own implementation adding a function `func main()` in the `main` package.

If you send a zip file, you have to provide your implementation in a file called `exec` (without extension `.go`!) placed in the top level of the zip file.

If you provide your own `main.main()`, the default `main` will not be generated.

An example named `standalone` is provided.
