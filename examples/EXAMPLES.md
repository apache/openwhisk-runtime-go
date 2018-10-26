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
# Examples


This is a collection of examples.
Tested on:

- dep version 0.5.0 (check the version)
- Go version 1.11.1
- GNU Make 3.81
- Linux Ubuntu 14.04
- Mac OSX 10.13
- Windows, with Git Bash, Docker for Windows, make from ezwinports, zip.exe

Each examples has a  Makefile with 4 targets:

- `make deploy` (or just make) deploys the action, precompiling it
- `make devel`  deploys the action in source format, leaving the compilation to the runtime
-  `make test` runs a simple test on the action; it should be already deployed
- `clean` removes intermediate files

Available examples:

- [Simple Golang action](golang-main-single) main is `main.Main`
- [Simple Golang action](golang-hello-single) main is `main.Hello`
- [Golang action with a subpackage](golang-main-package) main is `main.Main` invoking a `hello.Hello`
- [Golang action with a vendor folder](golang-main-vendor) main is `main.Main` using a dependency `github.com/rs/zerolog`
- [Golang action with a subpackage and vendor folder](golang-hello-vendor) main is `main.Hello` invoking a `hello.Hello` using a dependency `github.com/sirupsen/logrus`
- [Standalone Golang Action](golang-main-standalone) main is `main.main`, implements the ActionLoop directly
- [Simple Bash action](bash-hello) a simple bash script action implementing the ActionLoop directly



