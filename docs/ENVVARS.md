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

# Environment Variables

The proxy uses a few enviroment variables that affects its behavior.

`OW_COMPILER` points to the compiler script to use to compile actions.

`OW_SAVE_JAR` enables checking that an upload file is a jar (that is itself a zip file) and it will not expand it if there is a subdirectory names "META-INF" (so it is a jar file). Used to support uploading of Java jars.

`OW_WAIT_FOR_ACK` enables waiting for an acknowledgement in the actionloop protocol. It should be enabled in all the newer runtimes. Do not enable in existing runtimes as it would break existing actions built for that runtime.

`OW_EXECUTION_ENV` enables detection and verification of the compilation environent. The compiler is expected to create a file named `exec.env` in the same folder as the `exec` file to be run. If this variable is set, before starting an action, the init will check that the content of the `exec.env` starts with the value of the variable. The actual content of the `exec.env` can be actually a longer string. 

This variable is also propagated to the action as `__OW_EXECUTION_ENV`. Furthermore also the version of the proxy is propagated to the action as  `__OW_EXECUTION_VERSION`.




