#!/usr/bin/env python
# Licensed to the Apache Software Foundation (ASF) under one or more contributor
# license agreements; and to You under the Apache License, Version 2.0.
import os, sys
os.rename(sys.argv[2], sys.argv[3]+"/action")
with open(sys.argv[3]+"/exec", "w") as f:
    f.write("""#!/bin/bash
cd "$(dirname $0)"
export PYTHONPATH=$PWD/action
python action/exec.py
""")
