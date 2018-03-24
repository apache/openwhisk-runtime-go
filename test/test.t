  $ export T=$TESTDIR

  $ $T/bin/run.sh 
  {"error":"no action defined yet"}
  400

  $ $T/bin/post.sh $T/etc/empty.json
  {"ok":true}
  200

  $ $T/bin/init.sh $T/test.t
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

  $ $T/bin/init.sh $T/test.t
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