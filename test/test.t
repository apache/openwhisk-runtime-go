  $ export T=$TESTDIR
  $ $T/build.sh >/dev/null

  $ $T/bin/run.sh 
  {"error":"no action defined yet"}

  $ $T/bin/init.sh $T/test.t
  {"error":"cannot start action:* (glob)
  400

  $ $T/bin/init.sh $T/bin/empty
  {"error":"cannot start action: sent invalid action"}
  400

  $ $T/bin/init.sh $T/bin/hi
  {"error":"cannot start action: sent invalid action"}
  400

  $ $T/bin/run.sh 
  {"error":"no action defined yet"}
  400

  $ $T/bin/init.sh $T/bin/hello_message
  {"ok":true}
  200

  $ $T/bin/run.sh 
  {"message":"Hello, Mike!"}

  $ $T/bin/init.sh $T/bin/hello_greeting
  {"ok":true}
  200

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}

  $ $T/bin/init.sh $T/zip/hello_message.zip
  {"ok":true}

  $ $T/bin/run.sh 
  {"message":"Hello, Mike!"}
  
  $ $T/bin/init.sh $T/zip/hello_greeting.zip
  {"ok":true}
  200

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}

  $ $T/bin/init.sh $T/test.t
  {"error":"cannot start action:* (glob)
  400

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}

  $ $T/bin/init.sh $T/bin/empty
  {"error":"cannot start action: sent invalid action"}
  400

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}

  $ $T/bin/init.sh $T/bin/hi
  {"error":"cannot start action: sent invalid action"}
  400

  $ docker kill $(cat $T/last.id) >/dev/null
