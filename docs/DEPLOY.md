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

There are two images provided: the `actionloop` and the `actionloop-golang-v1.11` available. Each image accept different input in the deployment.

<a name="actionloop">

## Actionloop runtime

The runtime `actionlooop` accepts:

- single file executable
- a zip file containing an executables

If the input is a single file, it can be either a in ELF format for architecture AMD64 implementing the ActionLoop protocol.

It can also be a script, identified by the `#!` hash-bang path at the beginning. The default `actionloop` can execute `bash` shell scripts and can use the `jq` command to parse JSON files and the `curl` command to invoke other actions.

If the file is a zipped file, it must contain in the top level (*not* in a subdirectory) an file named `exec`. This file must be in the same format as a single binary, either a binary or a script.

<a name="golang">

## Golang runtime

The runtime `actionloop-golang-v1.11` accepts:

- executable binaries implementing the ActionLoop protocol as Linux ELF executable compiled for the AMD64 architecture (as the `actionloop` runtme)
- zip files containing a binary executable named `exec` in the top level, and it must be again a Linux ELF executable compiled for the AMD64 architecture
- a single file action that is not an executable binary will be interpreted as source code and it will be compiled in a binary as described in the document about [actions](ACTION.md)
- a zip file not containing in the top level a binary file `exec` will  be interpreted as a collection of zip files, and it will be compiled in a binary as described in the document about [actions](ACTION.md)

Please note in the separate the rules about the name of the main function (that defaults to `main.Main`), and the rules about how to overwrite the `main.main`.

## Using packages and vendor folder

When you deploy a zip file, you can:

- have all your functions in the `main` package
- have some functions placed in some packages, like `hello`
- have some third party dependencies you want to include in your sources

If all your functions are in the main package, just place all your sources in the top level of your zip file

### Use a package folder

If some functions belong to a package, like `hello/`, you need to be careful with the layout of your source. The layout supported is the following:

```
golang-main-package/
├── Makefile
└── src
    ├── hello
    │   ├── hello.go
    │   └── hello_test.go
    └── main.go
```

You need to use a `src` folder, place the sources that belongs to the main package in the `src` and place sources of your package in the `src/hello` folder.

Then you should import it your subpackage with `import "hello"`.
Note that this means if you want to compile locally you have to set your GOPATH to parent directory of your `src` packages. Check below for using [VcCode](#vscode) as an editor with this setup.

When you send the image you will have to zip the content

Check the example `golang-main-package` and the associated `Makefile` for an example including also how to deploy and precompile your sources.

### Using vendor folders

When you need to use third part libraries, the runtime does not download them from Internet. You have to provide them,  downloading and placing them using the `vendor` folder mechanism. We are going to show here how to use the vendor folder with the `dep` tool.

*NOTE* the `vendor` folder does not work at the top level, you have to use a `src` folder and a package folder to have also the vendor folder.

If you want for example use the library `github.com/sirupsen/logrus` to manage your logs (a widely used drop-in replacement for the standard `log` package), you have to include it in your source code *in a sub package*.

For example consider you have in the file `src/hello/hello.go` the import:

```
import "github.com/sirupsen/logrus"
```

To create a vendor folder, you need to

- install the [dep](https://github.com/golang/dep) tool
- cd to the `src/hello` folder (*not* the `src` folder)
- run `GOPATH=$PWD/../.. dep init` the first time (it will create 2 manifest files `Gopkg.lock` and `Gopkg.toml`) or `dep ensure` if you already have the manifest files.

The layout will be something like this:

```
golang-hello-vendor
├── Makefile
└── src
    ├── hello
    │   ├── Gopkg.lock
    │   ├── Gopkg.toml
    │   ├── hello.go
    │   ├── hello_test.go
    │   └── vendor
    │       ├── github.com/...
    │       └── golang.org/...
    └── hello.go
```

Check the example `golang-hello-vendor`.

Note you do not need to store the `vendor` folder in the version control system as it can be regenerated (only the manifest files), but you need to include the entire vendor folder when you deploy the action.

If you need to use vendor folder in the main package, you need to create a directory `main` and place all the source code that would normally go in the top level, in the `main` folder instead.  A vendor folder in the top level *does not work*.

<a name="vscode">

### Using VsCode

If you are using [VsCode[(https://code.visualstudio.com/) as your Go development environment with the [VsCode Go](https://marketplace.visualstudio.com/items?itemName=ms-vscode.Go) support, and you want to get rid of errors and have it working properly, you need to configure it to support the suggested:

- you need to have a `src` folder in your source
- you need either to open the `src` folder as the top level source or add it as a folder in the workspace (not just have it as a subfolder)
- you need to enable the option `go.inferGopath`

Using this option, the GOPATH will be set to the parent directory of your `src` folder and you will not have errors in your imports.

<a name="precompile"/>

## Precompiling Go Sources Offline

Compiling sources on the image can take some time when the images is initialized. You can speed up precompiling the sources using the image `actionloop-golang-v1.11` as an offline compiler. You need `docker` for doing that.

The images accepts a `-compile <main>` flag, and expects you provide sources in standard input. It will then compile them, emit the binary in standard output and errors in stderr. The output is always a zip file containing an executable.

If you have docker, you can do it this way:

If you have a single source maybe in file `main.go`, with a function named `Main` just do this:

`docker run openwhisk/actionloop-golang-v1.11 -compile main <main.go >main.zip`

If you have multiple sources in current directory, even with a subfolder with sources, you can compile it all with:

`zip -r - * | docker run openwhisk/actionloop-golang-v1.11 -compile main >main.zip`

The  generated executable is suitable to be deployed in OpenWhisk using just the generic `actionloop` runtime.

`wsk action create my/action main.zip -docker openwhisk/actionloop`

You can also use the full `actionloop-golang-v1.11` as runtime, it is only bigger.

Note that the output is always a zip file in  Linux AMD64 format so the executable can be run only inside a Docker Linux container.

<a name="knative"/>
# Knative support

Action precompilation can be performed in a Knative build, producing a precompiled action. You can then bundle the action in a single image.

If you set the ennvironmente variables `OW_AUTOINIT` to the image executable, the action is automatically initialized. 

If the action is a binary executable a `main` is not required, however you can also bundle a zip file containing sources. In this case you may need to specify also the main function, using the environment variable `OW_AUTOINIT_MAIN`.

In the folder `examples/knative` there is an example building an images with Tekton Pipelines.






