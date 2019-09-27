#!/bin/bash
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
kubectl -n sample get po -w
