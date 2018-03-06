  $ export T=$TESTDIR

  $ $T/bin/run.sh 
  {"error":"no action defined yet"}

  $ $T/bin/init.sh $T/test.t
  {"error":"invalid action:* (glob)

  $ $T/bin/init.sh $T/bin/empty
  {"error":"cannot start action: sent invalid action"}

  $ $T/bin/init.sh $T/bin/hi
  {"error":"cannot start action: sent invalid action"}

  $ $T/bin/run.sh 
  {"error":"no action defined yet"}

  $ $T/bin/init.sh $T/bin/hello_message
  OK

  $ $T/bin/run.sh 
  {"message":"Hello, Mike!"}

  $ $T/bin/init.sh $T/bin/hello_greeting
  OK

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}

  $ $T/bin/init.sh $T/zip/hello_message.zip
  OK

  $ $T/bin/run.sh 
  {"message":"Hello, Mike!"}

  $ $T/bin/init.sh $T/zip/hello_greeting.zip
  OK

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}

  $ $T/bin/init.sh $T/test.t
  {"error":"invalid action:* (glob)

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}

  $ $T/bin/init.sh $T/bin/empty
  {"error":"cannot start action: sent invalid action"}

  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}

  $ $T/bin/init.sh $T/bin/hi
  {"error":"cannot start action: sent invalid action"}
