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

This is an optional command line test kit.

It has been superseded by a go test suite and scalatests,
but it is still around for debugging.


TO use it:
- you need to install cram to use it and build images
- you can run them with `cram test_actionloop.t`
- and `cram test_actionloop-golang`

also you can start directly the binary without the images with

- `./start.sh`
- or `COMPILER=../common/gobuild.sh ./start.sh`

If you start them,  images won't be started by the test.



