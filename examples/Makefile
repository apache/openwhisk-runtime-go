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

# prefix
OW_USER?=openwhisk
OW_VER?=v1.16:nightly

.PHONY: test
test: all
	cat */test.out | grep 'Hello' | sort >test.out
	diff test.out test.expected

.PHONY: all
all: bash standalone single-main single-hello package-main module-main

.PHONY: bash
bash:
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean deploy test

.PHONY: standalone
standalone:
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean devel test
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean deploy test

.PHONY: single-main
single-main:
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean devel test
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean deploy test

.PHONY: single-hello
single-hello:
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean devel test
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean deploy test

.PHONY: package-main
package-main:
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean devel test
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean deploy test

.PHONY: module-main
module-main:
	# this exceeds 3 seconds for initialization
	#cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean devel test
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean deploy test

.PHONY: aws-main
aws-main:
	# this exceeds 3 seconds for initialization
	#cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean devel test
	cd $@ && OW_USER=$(OW_USER) OW_VER=$(OW_VER) make clean deploy test
