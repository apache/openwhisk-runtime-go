#!/bin/bash
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

cd "$(dirname $0)"

function build {
   test -e exec && rm exec
   cp $1.src $1.go
   GOPATH=$PWD go build -a -o exec $1.go
   rm $1.go
}

function build_main {
   test -e exec && rm exec
   cp ../../common/gobuild.py.launcher.go $1.go
   cat $1.src >>$1.go
   go build -a -o exec $1.go
   rm $1.go
}


build hi
zip -q hi.zip exec
cp exec hi

build_main hello_message
zip -q hello_message.zip exec
cp exec hello_message

build_main hello_greeting
zip -q hello_greeting.zip exec
cp exec hello_greeting

test -e hello.zip && rm hello.zip
cd src
zip -q -r ../hello.zip main.go hello
cd ..

test -e sample.jar && rm sample.jar
cd jar ; zip -q -r ../sample.jar * ; cd ..

build exec
test -e exec.zip && rm exec.zip
zip -q -r exec.zip exec etc dir
echo exec/env >helloack/exec.env
zip -j helloack.zip helloack/*

python3 -m venv venv
