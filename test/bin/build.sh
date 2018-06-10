#!/bin/sh
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
FILE=${1:?go file}
EXEC=${2:-main}
OUT=$(basename $FILE)
BIN=${OUT%%.go}
ZIP=${BIN}.zip
rm bin/$BIN
go build -i -o bin/$BIN $FILE
GOOS=linux GOARCH=amd64 go build -o $EXEC $FILE
rm zip/$ZIP
zip zip/$ZIP $EXEC
rm $EXEC
echo "built $EXEC bin/$BIN zip/$ZIP"
