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

This is an (optional) command line test kit.

There is now also go test suite and scalatests,
but we still use it for debugging.

- It works on Linux Ubuntu only (on OSX there are some differences in the CLI).
- You need to build the images with `dockerDist`.
- You need to install cram .
- Test action loop: `cram test_actionloop.t`
- Test golang `cram test_actionloop-golang.t`

Also you can start directly the executable, using `start.sh`  without building the images.
So you can debug outside of Docker.
If you start the executable  images won't be started by the test.

- `unset COMPILER ; ./start.sh` for action loop, then  `cram test_actionloop.t`
- or `COMPILER=../common/gobuild.sh ./start.sh` then `cram test_actionlooop-golang.t`




