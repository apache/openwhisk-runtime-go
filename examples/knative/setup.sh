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
USER=${1:?user}
PASS=${2:?pass}
HOST=${3:?host}

kubectl apply -f - <<EOF
---
apiVersion: v1
kind: Namespace
metadata:
  name: sample
---
apiVersion: v1
kind: Secret
metadata:
  name: docker-user-pass
  namespace: sample
  annotations:
    tekton.dev/docker-0: https://$HOST
type: kubernetes.io/basic-auth
stringData:
  username: $USER
  password: $PASS
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: docker-service-account
  namespace: sample
secrets:
  - name: docker-user-pass
EOF
kubectl -n sample apply -f 0-buildkitd.yaml
kubectl -n sample apply -f 1-task.yaml

