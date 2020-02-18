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
IMAGES?=openwhisk

all: .PHONY

.PHONY: bash-hello golang-hello-single golang-main-single golang-main-standalone golang-main-package golang-main-vendor golang-hello-vendor

bash-hello:
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-main-single:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-hello-single:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-main-standalone:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-main-package:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-main-vendor:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-hello-vendor:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test


