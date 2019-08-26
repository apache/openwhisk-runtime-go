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
# ActionLoop v1.0.1
- embedded file type detection
- now showing the commend
- librdkafka in golang image
- showing version numbuer with -debug

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

# Deincubation
- Removed the -incubation
- Now all runtimes use source release so no more actionloop-v2, renamed to actionloop-base
- upgraded to go 1.12.9 and 1.11.13
