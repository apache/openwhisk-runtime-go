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
WSK?=wsk
OW_USER?=openwhisk
OW_RUNTIME?=$(OW_USER)/actionloop-base:nightly
NAME=bash
PACKAGE=test
SRC=hello.sh

deploy: package.done
	$(WSK) action update $(PACKAGE)/$(NAME) $(SRC) --docker $(OW_RUNTIME)

devel: deploy

test: test.json
	$(WSK) action invoke $(PACKAGE)/$(NAME) -r | tee -a test.out
	$(WSK) action invoke $(PACKAGE)/$(NAME) -P test.json -r | tee -a test.out

clean:
	-$(WSK) action delete $(PACKAGE)/$(NAME)
	-rm package.done test.json
	-rm test.out

package.done:
	$(WSK) package update $(PACKAGE)
	touch package.done

test.json:
	echo '{ "name": "Mike" }' >test.json

.PHONY: test devel deploy clean
