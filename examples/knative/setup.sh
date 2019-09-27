#!/bin/bash
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

