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
# Knative Support

This is an example to build an action with Tekton Pipelines.

As a prerequisite you need on the path `kubectl` configured to point to a Kubernetes cluster with [Tekton Pipelines](https://github.com/tektoncd/pipeline/blob/master/docs/install.md) already installed.

You also need a  docker registry. For example you can get a free account for public images on Docker Hub.

You need to put your sources in a git repository, like GitHub. As an example you can use `https://github.com/sciabarracom/hellogo`. This sample expects the source of an action in a file named `src` and a Dockerfile to add the resoult of the compilation named `exec.zip` to the image of the runtime.

To use this example with GitHub and Docker Hub, first initialize the Tekton Build with:

`./setup.sh <docker-user> <docker-password> index.docker.io`

you can then build with

`./build.sh <git-source> <docker-image>`

Using the example in GitHub to build an image in DockerHub you can use (change `actionloop` to your DockerHub user)

`./build.sh https://github.com/sciabarracom/hellogo docker.io/actionloop/hellogo`

You can watch the build status with `kubectl -n sample get po -w` until it completes.
