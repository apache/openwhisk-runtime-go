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
FROM golang:1.15.0
RUN echo "deb http://deb.debian.org/debian stretch-backports main contrib non-free" \
     >>/etc/apt/sources.list &&\
    apt-get update && apt-get install -y \
     curl \
     jq \
     git \
     vim &&\
    apt-get -y install \
     librdkafka1=0.11.6-1~bpo9+1 \
     librdkafka++1=0.11.6-1~bpo9+1 && \
    apt-get -y install \
     librdkafka-dev=0.11.6-1~bpo9+1 && \
    rm -rf /var/lib/apt/lists/* && \
    go get -u github.com/go-delve/delve/cmd/dlv
RUN mkdir /action
WORKDIR /action
ADD proxy /bin/proxy
ADD bin/compile /bin/compile
ADD lib/launcher.go /lib/launcher.go
ENV OW_COMPILER=/bin/compile
ENV OW_LOG_INIT_ERROR=1
ENV OW_WAIT_FOR_ACK=1
ENV OW_EXECUTION_ENV=openwhisk/action-golang-v1.15
ENTRYPOINT [ "/bin/proxy" ]
