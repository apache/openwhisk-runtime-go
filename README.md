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
# Apache OpenWhisk Runtimes for Go

[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0)
[![Continuous Integration](https://github.com/apache/openwhisk-runtime-go/actions/workflows/ci.yaml/badge.svg)](https://github.com/apache/openwhisk-runtime-go/actions/workflows/ci.yaml)
[![Join Slack](https://img.shields.io/badge/join-slack-9B69A0.svg)](http://slack.openwhisk.org/)

This repository contains both the OpenWhisk runtime for Golang Actions, as well as a runtime for Generic executables.

- If you are in a hurry, check the [examples](examples/EXAMPLES.md)
- Writing Actions for the runtime in [Golang](docs/ACTION.md#golang)
- How to deploy your [Golang](docs/DEPLOY.md#golang) sources
- Precompiling [Golang](docs/DEPLOY.md#precompile) actions
- How to use VSCode to write [Golang](docs/DEPLOY.md#vscode) actions
- How to [Build](docs/BUILD.md#building) the runtime, with development notes

## Actionloop runtime

### Using the Go runtime for Generic executables

- Writing [Generic](docs/ACTION.md#generic) actions, in bash or as a generic linux binary
- Deployment for [Generic](docs/DEPLOY.md#generic) actions
- The [ActionLoop](docs/ACTION.md#actionloop) protocol for generic actions
- Environment [Variables](docs/ENVVARS.md) to configure the proxy

# Change Log

[CHANGES.md](CHANGES.md)

# License
[Apache 2.0](LICENSE.txt)
