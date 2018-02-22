#!/bin/bash
rm -Rvf Packages.resolved .build
rm exec
docker run -v $PWD:/root ibmcom/swift-ubuntu:4.0.3 bash -c "cd /root && swift build && cp .build/debug/Hello exec"
