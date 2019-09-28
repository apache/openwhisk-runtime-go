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
USR?=opennwhisk
BASE=actionloop-golang-v1.12
IMG=docker.io/$(USR)/$(BASE)
OUT=hello-go

.PHONY: push
push:
	docker tag $(BASE) $(IMG)
	docker login
	docker push $(IMG)

.PHONY: test-compile
test-compile:
	docker run -i $(IMG) \
	-compile main \
	 <../golang-main-single/main.go \
	 >exec.zip
	 unzip -l exec.zip

.PHONY: test-build
test-build:
	mkdir workspace 
	git clone https://gist.github.com/sciabarracom/a84ce9baa7d8eb7a769f80c61a618694 workspace/source
	docker run -i --entrypoint=/bin/bash \
	-v $(PWD)/workspace:/workspace \
	$(IMG) -c "/bin/proxy -compile main </workspace/source/src >/workspace/source/exec.zip"
	docker build -t $(OUT) workspace/source 

.PHONY: enter
enter:
	docker run \
	  -v $(PWD)/workspace:/workspace \
	  --entrypoint=/bin/bash \
	  -ti $(IMG)



