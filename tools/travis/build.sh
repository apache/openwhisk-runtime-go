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

set -ex

# Build script for Travis-CI.

SCRIPTDIR=$(cd "$(dirname "$0")" && pwd)
ROOTDIR=$(cd "$SCRIPTDIR/../.." && pwd)
WHISKDIR=$(cd "$ROOTDIR/../openwhisk" && pwd)

export OPENWHISK_HOME=$WHISKDIR

IMAGE_PREFIX="testing"

# Build OpenWhisk
cd "$WHISKDIR"

#pull down images
docker pull openwhisk/controller
docker tag openwhisk/controller ${IMAGE_PREFIX}/controller
docker pull openwhisk/invoker
docker tag openwhisk/invoker ${IMAGE_PREFIX}/invoker

./gradlew --console=plain \
:common:scala:install \
:core:controller:install \
:core:invoker:install \
:tests:install

# For pull requests, force a local-only build
if [ "$TRAVIS_PULL_REQUEST" != "false" ]; then
  export docker_local_json='{"amd64":null}'
fi

builds=( ":core:golangAction:dockerBuildImage" )

# Build runtime
echo "---------------------------------------------------------------------------------------"
echo " Building " "${builds[@]}"
echo "---------------------------------------------------------------------------------------"

cd "$ROOTDIR"
./gradlew --console=plain "${builds[@]}" -PdockerImagePrefix=${IMAGE_PREFIX}

echo "---------------------------------------------------------------------------------------"

