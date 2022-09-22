#!/usr/bin/env python
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

from __future__ import print_function

import os
import json

from action.main import main
inp = os.fdopen(0, "rb")
out = os.fdopen(3, "wb")
while True:
    while True:
        line = inp.readline()
        args = json.loads(line)
        payload = {}
        if "value" in args:
            payload = args["value"]
        res = main(payload)
        out.write(json.dumps(res, ensure_ascii=False).encode('utf-8'))
        out.write(b"\n")
        out.flush()
