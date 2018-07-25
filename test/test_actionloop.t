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

  $ export T=$TESTDIR
  $ $T/start.sh actionloop 2>$T/err.log >$T/out.log 

Default: no action defined 

  $ $T/bin/post.sh $T/etc/empty.json
  {"ok":true}
  200

  $ $T/bin/run.sh 
  {"error":"no action defined yet"}
  400

  $ $T/bin/init.sh $T/etc/hello.go
  {"error":"cannot start action: command exited"}
  400

  $ $T/bin/init.sh $T/bin/empty
  {"error":"cannot start action: command exited"}
  400

  $ $T/bin/init.sh $T/bin/hi
  {"error":"cannot start action: command exited"}
  400

  $ $T/bin/run.sh 
  {"error":"no action defined yet"}
  400

Sending some messages 

  $ $T/bin/init.sh $T/bin/hello_message
  {"ok":true}
  200

  $ $T/bin/run.sh 
  {"message":"Hello, Mike!"}
  200

  $ $T/bin/init.sh $T/bin/hello_greeting
  {"ok":true}
  200

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}
  200

  $ $T/bin/init.sh $T/zip/hello_message.zip
  {"ok":true}
  200

  $ $T/bin/run.sh 
  {"message":"Hello, Mike!"}
  200

  $ $T/bin/init.sh $T/zip/hello_greeting.zip
  {"ok":true}
  200

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}
  200

Invalid  script

  $ $T/bin/init.sh $T/test_actionloop.t
  {"error":"cannot start action: command exited"}
  400

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}
  200

  $ $T/bin/init.sh $T/bin/empty
  {"error":"cannot start action: command exited"}
  400

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}
  200

  $ $T/bin/init.sh $T/bin/hi
  {"error":"cannot start action: command exited"}
  400

A shell script

  $ $T/bin/init.sh $T/etc/hello.sh
  {"ok":true}
  200

  $ $T/bin/run.sh
  {"hello":"Hello, Mike"}
  200

  $ $T/bin/run.sh '{"name": ""}'
  {"error":"command exited"}
  400

  $ $T/bin/run.sh
  {"error":"no action defined yet"}
  400


Test with a non-main executable

  $ $T/bin/init.sh $T/zip/hello_message1.zip
  {"error":"cannot start action: command exited"}
  400

  $ $T/bin/init.sh $T/zip/hello_message1.zip hello_message
  {"error":"cannot start action: command exited"}
  400

  $ $T/stop.sh 
