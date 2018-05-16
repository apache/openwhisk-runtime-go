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
if netstat -tlnp 2>/dev/null | grep :8080 >/dev/null
then exit 0
fi
cd "$(dirname $0)"
bin/build.sh src/hello_greeting.go
bin/build.sh src/hello_message.go
bin/build.sh src/empty.go
bin/build.sh src/hi.go
zip -j zip/hello_message1.zip bin/hello_message
rm -Rvf action
mkdir action
if test -n "$1"
then docker run --name=goproxy -p 8080:8080 -d "$@"
else go run ../main/proxy.go -debug
fi
