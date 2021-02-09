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
# 1.17.0
- go 1.15 runtime upgraded to 1.15.7
- go 1.13 runtime upgraded to 1.13.15
- add 'apt-get upgrade' to the image build of go 1.15 and go 1.13 to get latest security fixes during each build, for the case the base images are not updated frequently

# 1.16.0
- added go 1.13 and 1.15 with Go modules
- removed support for go1.11 and go1.12
- updated examples
- add 'apt-get upgrade' to the image build to get latest security fixes during each build, for the case the base images are not updated frequently
- added OW_WAIT_FOR_ACK such at if true, the proxy waits for an acknowledgement from the action on startup
- added OW_EXECUTION_ENV to validate the execution environment before starting an action
- write compilation logs to standard out
# 1.15.0
- added OW_ACTION_VERSION to action environment (PR#113)
- propagate API_HOST from parent to child process (PR#115)

# 1.14.0
- Removed the -incubation
- Now all runtimes use source release so no more actionloop-v2, renamed to actionloop-base
- upgraded to go 1.12.9 and 1.11.13

# Actionloop v2
Versioning
- renamed actionloop docker image to actionloop-v2
Docker Images Support
- static build of the executable docker image, so actionloop can be used also in alpine images
ActionLoop for Scripting Languages
- any script starting with '#!' is recognized as executable
- now the -compile will zip the entire directory of the `bin` directory after compilation
- if you upload a folder `src/exec` the entire directory is moved to `bin`, including other uploaded files
- Support for Go 1.12.4
- Support for jar not expanded for Java when set OW_SAVE_JAR
- You can initalize multiple times when debugging
- Removed gogradle plugin, now building directly with go

# ActionLoop v1.0.1
- embedded file type detection
- now showing the commend
- librdkafka in golang image
- showing version numbuer with -debug
