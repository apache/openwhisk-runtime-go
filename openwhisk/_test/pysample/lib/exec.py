from __future__ import print_function
import os, json
from action.main import main
inp = os.fdopen(0, "rb")
out = os.fdopen(3, "wb")
while True:
    while True:
       line = inp.readline()
       args = json.loads(line)
       payload = {} 
       if "value" in args:
           payload = args["value"]
       res = main(payload)
       out.write(json.dumps(res, ensure_ascii=False).encode('utf-8'))
       out.write("\n")
       out.flush()

