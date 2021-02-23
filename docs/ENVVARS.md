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

## Environment variables that control the behavior of the proxy

The following variables are usually set in the Dockerfile

`OW_COMPILER` points to the compiler script to use to compile actions.

`OW_SAVE_JAR` enables checking that an uploaded file is a jar (that is itself a zip file) and it will not expand it if there is a subdirectory named "META-INF" (so it is a jar file). Used to support uploading of Java jars.

`OW_WAIT_FOR_ACK` enables waiting for an acknowledgment in the action loop protocol. It should be enabled in all the newer runtimes. Do not enable in existing runtimes as it would break existing actions built for that runtime.

`OW_EXECUTION_ENV` enables detection and verification of the compilation environment. The compiler is expected to create a file named `exec.env` in the same folder as the `exec` file to be run. If this variable is set, before starting an action, the initialization will check that the content of the `exec.env`, trimmed of spaces and new lines, is the same, to ensure an action is executed in the right execution environment.

`OW_LOG_INIT_ERROR` enables logging of compilation error; the default behavior is to return errors in the result from initialization.

## Environment variables propagated to actions and to the compilation script

The proxy itself sets the following environment variables:

`__OW_EXECUTION_ENV` is the same value that the proxy receives as `OW_EXECUTION_ENV`

`__OW_WAIT_FOR_ACK` is set if the proxy has the variable `OW_WAIT_FOR_ACK` set.

Any other environment variables set in the Dockerfile that start with `__OW_` are propagated to the proxy and can override the values set by the proxy.

Furthermore, actions receive their own environment variables and such values override the variables set from the proxy and in the environment.
