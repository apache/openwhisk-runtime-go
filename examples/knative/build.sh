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
SRC=${1:?git source}
TGT=${2:?image target}
ID=$(date +%s)
kubectl -n sample apply -f - <<EOF
apiVersion: tekton.dev/v1alpha1
kind: TaskRun
metadata:
  name: knative-build-$ID
spec:
  serviceAccount: docker-service-account
  taskRef:
    name: buildkit
  inputs:
    resources:
    - name: source
      resourceSpec:
        type: git
        params:
        - name: url
          value: $SRC
  outputs:
    resources:
    - name: image
      resourceSpec:
        type: image
        params:
        - name: url
          value: $TGT
EOF
echo "check status with: kubectl -n sample get po -w"
