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

# Deployment options

There are a few images provided: the `actionloop-base` and the `action-golang-v1.19` and `action-golang-v1.20` are available. Each image accepts different input in the deployment.

<a name="actionloop">

## Actionloop runtime

The runtime `actionlooop-base` accepts:

- single file executable
- a zip file containing an executables

If the input is a single file, it can be either a in ELF format for architecture AMD64 implementing the ActionLoop protocol.

It can also be a script, identified by the `#!` hash-bang path at the beginning. The default `actionloop-base` can execute `bash` shell scripts and can use the `jq` command to parse JSON files and the `curl` command to invoke other actions.

If the file is a zipped file, it must contain in the top level (*not* in a subdirectory) an file named `exec`. This file must be in the same format as a single binary, either a binary or a script.

<a name="golang">

## Golang runtime

The runtime `action-golang-v1.N` accepts:

- executable binaries implementing the ActionLoop protocol as Linux ELF executable compiled for the AMD64 architecture (as the `actionloop-base` runtime)
- zip files containing a binary executable named `exec` in the top level, and it must be again a Linux ELF executable compiled for the AMD64 architecture
- a single file action that is not an executable binary will be interpreted as source code and it will be compiled in a binary as described in the document about [actions](ACTION.md)
- a zip file not containing in the top level a binary file `exec` will  be interpreted as a collection of zip files, and it will be compiled in a binary as described in the document about [actions](ACTION.md)

Please note in the separate the rules about the name of the main function (that defaults to `main.Main`), and the rules about how to overwrite the `main.main`.

## Using packages and modules

When you deploy a zip file, you can:

- have all your functions in the `main` package
- have some functions placed in some packages, like `hello`
- have some third party dependencies you want to include in your sources

You can manage those dependencies using appropriate `go.mod` files using relative and absolute references.

For example you can use a local package `hello` with:

```
replace hello => ./hello
```

Check the example: `package-main` and `module-main` and look for the format of the `go.mod` files.

<a name="precompile"/>

## Precompiling Go Sources Offline

Compiling sources on the image can take some time when the images is initialized. You can speed up precompiling the sources using the image `action-golang-v1.N` as an offline compiler. You need `docker` for doing that.

The images accepts a `-compile <main>` flag, and expects you provide sources in standard input. It will then compile them, emit the binary in standard output and errors in stderr. The output is always a zip file containing an executable.

If you have docker, you can do it this way:

If you have a single source maybe in file `main.go`, with a function named `Main` just do this:

`docker run openwhisk/action-golang-v1.N -compile main <main.go >main.zip`

If you have multiple sources in current directory, even with a subfolder with sources, you can compile it all with:

`zip -r - * | docker run openwhisk/action-golang-v1.N -compile main >main.zip`

You can then execute the code. Note you have to use the same runtime you used to build the image.

Note that the output is always a zip file in  Linux AMD64 format so the executable can be run only inside a Docker Linux container.
