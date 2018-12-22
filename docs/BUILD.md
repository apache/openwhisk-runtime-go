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
# Developers Guide for Runtime itself

<a name="building"/>

## How to Build and run tests

You need a Linux or an OSX environment, with Java and Docker installed to build the sources.

Prerequisites for running build and tests with gradle:

-   docker
-   jdk

To compile go proxy *in your local architecture* for docker:

```
./gradlew goBuild
```

To build the docker images (this will also compile the go proxy):

```
./gradlew distDocker
```

This will build the images:

*   `actionloop-golang-v1.11`: an image supporting  Go sources
*   `actionloop`: the base image, supporting generic executables and shell script

The `actionloop` image can be used for supporting other compiled programming languages as long as they implement a `compile` script and the *action loop* protocol described below.

To run tests:

```
./gradlew test --info
```
<a name="development"/>

# Local Development

If you want to develop the proxy and run tests natively, you can do it on Linux or OSX. Development has been tested on Ubuntu Linux (14.04) and OSX 10.13. Probably other distributions work, maybe even Windows with WSL, but since it is not tested YMMMV.

You need to install, of course [go 1.11.x](https://golang.org/doc/install)

Then you need a set of utilities used in tests:

-   bc
-   zip
-   realpath

Linux: `apt-get install bc zip realpath`
OSX: `brew install zip coreutils`

**NOTE**: Because tests build and cache some binary files, perform a `git clean -fx` and **do not share folders between linux and osx** because binaries are in different format...

You can also run the tests in go, without using `gradle` with

```
cd openwhisk
go test
```
