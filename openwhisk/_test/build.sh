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
   if test -e $1
   then return
   fi
   cp $1.src $1.go
   go build -a -o $1 $1.go
   rm $1.go
}

function zipit {
    if test -e $1
    then return
    fi
    mkdir $$
    cp $2 $$/$3
    zip -q -j $1 $$/$3
    rm -rf $$
}

go get github.com/apache/incubator-openwhisk-runtime-go/openwhisk

build exec
rm exec.zip
zip -q -r exec.zip exec etc dir

build hi
zipit hi.zip hi main

build hello_message
zipit hello_message.zip hello_message main
zipit hello_message1.zip hello_message message

build hello_greeting
zipit hello_greeting.zip hello_greeting main
zipit hello_greeting1.zip hello_greeting greeting

