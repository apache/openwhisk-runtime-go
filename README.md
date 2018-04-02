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

# Apache OpenWhisk runtimes for Go

[![Build Status](https://travis-ci.org/apache/incubator-openwhisk-runtime-nodejs.svg?branch=master)](https://travis-ci.org/apache/incubator-openwhisk-runtime-nodejs)

## Give it a try today

Create a zip action with a `exec` in the root of the zip

```
echo \
'#!/bin/bash
echo "{\"message\":\"Hello World\"}"' > exec
chmod +x exec
zip myAction.zip exec
```

Create the action using the docker image for the runtime

```
wsk action update myAction myAction.zip --docker openwhisk/goproxy:latest
```

This works on any deployment of Apache OpenWhisk

## To use on a deployment that contains the runtime as a kind

Create action using `--native`

```
wsk action update myAction myAction.zip --native
```

## Building and hacking

This runtime has been converted to a multi-architecture build.  For details on
building it, see
[here](https://github.com/apache/incubator-openwhisk/blob/master/docs/runtimes-building.md)

# License

[Apache 2.0](LICENSE.txt)
