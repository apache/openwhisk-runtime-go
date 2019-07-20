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

IMG?=whisk/actionloop-golang-v1.11
IMG2?=whisk/actionloop

all: golang bash

test.lua:
	echo 'wrk.method = "POST"'>test.lua
	echo "wrk.body = '{\"value\":{\"name\":\"Mike\"}}'">>test.lua
	echo 'wrk.headers["Content-Type"] = "application/json"'>>test.lua

golang: test.lua
	docker run -d --name under-test --rm -p 8080:8080 $(IMG)
	bash init.sh main.go
	wrk -t1 -c1 -stest.lua http://localhost:8080/run
	docker kill under-test

bash: test.lua
	docker run -d --name under-test --rm -p 8080:8080 $(IMG2)
	bash init.sh main.sh
	wrk -t1 -c1 -stest.lua http://localhost:8080/run
	docker kill under-test

.PHONY: all golang bash
