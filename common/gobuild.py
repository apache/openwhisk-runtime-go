#!/usr/bin/python
"""Golang Action Compiler
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
"""

import os
import re
import sys
import codecs
import subprocess

def sources(launcher, source_dir, main):

    func = main.capitalize()
    has_main = None

    # copy the exec to exec.go
    # also check if it has a main in it
    src = "%s/exec" % source_dir
    dst = "%s/exec__.go" % source_dir
    if os.path.isfile(src):
        with codecs.open(src, 'r', 'utf-8') as s:
            with codecs.open(dst, 'w', 'utf-8') as d:
                body = s.read()
                has_main = re.match(r".*package\s+main\W.*func\s+main\s*\(\s*\)", body, flags=re.DOTALL)
                d.write(body)

    # copy the launcher fixing the main
    if not has_main:
        dst = "%s/main__.go" % source_dir
        if os.path.isdir("%s/main" % source_dir):
            dst = "%s/main/main__.go" % source_dir
        with codecs.open(dst, 'w', 'utf-8') as d:
            with codecs.open(launcher, 'r', 'utf-8') as e:
                code = e.read()
                code = code.replace("Main", func)
                d.write(code)

def build(parent, source_dir, target):
    # compile...
    env = {
      "PATH": os.environ["PATH"],
      "GOPATH": os.path.abspath(parent),
      "GOCACHE": "/tmp",
      "GO111MODULE": "off"
    }
    if os.path.isdir("%s/main" % source_dir):
        source_dir += "/main"
    p = subprocess.Popen(
        ["go", "build", "-i", "-ldflags=-s -w",  "-o", target],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        cwd=source_dir,
        env=env)
    (o, e) = p.communicate()

    # stdout/stderr may be either text or bytes, depending on Python
    # version, so if bytes, decode to text. Note that in Python 2
    # a string will match both types; so also skip decoding in that case
    if isinstance(o, bytes) and not isinstance(o, str):
        o = o.decode('utf-8')
    if isinstance(e, bytes) and not isinstance(e, str):
        e = e.decode('utf-8')

    # remove the comments mentioning the folder in order to normalize output
    o = re.sub(r"# .*\n", "", o, flags=re.MULTILINE)
    e = re.sub(r"# .*\n", "", e, flags=re.MULTILINE)

    if o:
        sys.stdout.write(o)
        sys.stdout.flush()

    if e:
        sys.stderr.write(e)
        sys.stderr.flush()

def main(argv):
    if len(argv) < 4:
        print("usage: <main-file> <source-dir> <target-dir>")
        sys.exit(1)

    main = argv[1]
    source_dir = argv[2]
    target_dir = argv[3]

    parent = os.path.dirname(os.path.abspath(source_dir))
    target = os.path.abspath("%s/exec" % target_dir)

    sources(argv[0]+".launcher.go", source_dir, main)
    build(parent, source_dir, target)

if __name__ == '__main__':
    main(sys.argv)
