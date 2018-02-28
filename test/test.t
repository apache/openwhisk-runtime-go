  $ export T=$TESTDIR
  $ $T/bin/run.sh 
  {"error":"no action defined yet"}
  $ $T/bin/init.sh $T/bin/hello_greeting
  OK
  $ $T/bin/run.sh 
  {"greetings":"Hello, Mike"}
  $ $T/bin/init.sh $T/zip/hello_message.zip
  OK
  $ $T/bin/run.sh 
  {"message":"Hello, Mike!"}
  $ $T/bin/init.sh $T/test.t 
  {"error":"cannot write the file: no file"}
  $ $T/bin/run.sh 
  {"message":"Hello, Mike!"}

