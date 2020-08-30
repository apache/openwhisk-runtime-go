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

- Go version 1.14.6
- GNU Make 3.81
- Linux Ubuntu 18.04
- Mac OSX 10.15.6
- Windows 10 with WSL 2

Each examples has a  Makefile with 4 targets:

- `make deploy` (or just make) deploys the action, precompiling it
- `make devel`  deploys the action in source format, leaving the compilation to the runtime
- `make test` runs a simple test on the action; it should be already deployed
- `clean` removes intermediate files

Available examples:

- [Simple Golang action](single-main) main is `main.Main`
- [Simple Golang action](single-hello) main is `main.Hello`
- [Golang action with a package](package-main) main is `main.Main` invoking a `hello.Hello` and a test
- [Golang action with a module](module-main) main is `main.Main` using a dependency `github.com/rs/zerolog`
- [Standalone Golang Action](standalone) implements the ActionLoop directly in go
- [Simple Bash action](bash) action implementing the ActionLoop directly
