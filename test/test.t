  $ export T=$TESTDIR

no action

  $ $T/bin/run.sh 
  {"error":"no action defined yet"}

  $ $T/bin/init.sh $T/test.t
  {"error":"invalid action:* (glob)

  $ $T/bin/init.sh $T/bin/empty
  {"error":"cannot start action: cannot start action, deleted"}

  $ $T/bin/init.sh $T/bin/hi
  {"error":"cannot start action: cannot start action, deleted"}

  $ $T/bin/run.sh 
  {"error":"no action defined yet"}

  $ $T/bin/init.sh $T/bin/hello_message
  OK

  $ $T/bin/run.sh 
  {"error":"no action defined yet"}
