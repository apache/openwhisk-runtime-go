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

# Do not fix the patch level for golang:1.20 to automatically get security fixes.
FROM golang:1.20-bullseye

RUN echo "deb http://deb.debian.org/debian buster-backports main contrib non-free" \
     >>/etc/apt/sources.list &&\
    echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections &&\
    apt-get update &&\
    # Upgrade installed packages to get latest security fixes if the base image does not contain them already.
    apt-get upgrade -y --no-install-recommends &&\
    apt-get install -y apt-utils &&\
    apt-get install -y \
     curl \
     jq \
     git \
     zip \
     vim && \
     apt-get -y install \
     librdkafka1 \
     librdkafka++1 &&\
    apt-get -y install \
     librdkafka-dev &&\
    # Cleanup apt data, we do not need them later on.
    apt-get clean && rm -rf /var/lib/apt/lists/* &&\
    go install github.com/go-delve/delve/cmd/dlv@latest &&\
    mkdir /action
#make python 3 react as python
RUN ln -s /usr/bin/python3 /usr/bin/python

WORKDIR /action
ADD proxy /bin/proxy
ADD bin/compile /bin/compile
ADD lib/launcher.go /lib/launcher.go
ENV OW_COMPILER=/bin/compile
ENV OW_LOG_INIT_ERROR=1
ENV OW_WAIT_FOR_ACK=1
ENV OW_EXECUTION_ENV=openwhisk/action-golang-v1.20
ENTRYPOINT [ "/bin/proxy" ]
