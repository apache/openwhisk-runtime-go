#!/bin/bash
BASE="$(dirname $0)"
INVOKE=../../incubator-openwhisk-devtools/docker-compose/openwhisk-src/tools/actionProxy/invoke.py
OP=${1:?operation}
OPTION=${2:?option}
NAME=${3:?name}
FILE=${4}
PAYLOAD=$5

if test "$OP" == "action"
then
  case "$OPTION" in 
    update)
        python $BASE/$INVOKE init $FILE
    ;;
    invoke)
        python $BASE/$INVOKE run $PAYLOAD
    ;;
    *)
        echo ignored
    ;;
 esac
else echo "ignored"
fi
