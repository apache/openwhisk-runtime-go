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
FILE=${1:?file}
MAIN=${2:-}
JSON=/tmp/json$$
if test -n "$MAIN"
then MAIN=",\"main\":\"$MAIN\""
fi
if file -i $FILE | grep text/ >/dev/null
then echo '{"value":{"code":' $(cat $FILE | jq -R -s .) $MAIN '}}' >$JSON
else echo '{"value":{"binary":true,"code":"'$(base64 $FILE | tr -d '\n')'"}}' >$JSON
fi
curl -H "Content-Type: application/json" -XPOST -w "%{http_code}\n" http://${HOST:-localhost}:${PORT:-8080}/init -d @"$JSON" 2>/dev/null
#echo $JSON
#echo $MAIN
rm $JSON 2>/dev/null
