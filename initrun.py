import requests
import sys
import json
import six

def readfile(name):
    fp = open(name, "r")
    data = fp.read()
    fp.close()
    return data

#filename="hello.py"
def init(filename):
    body = json.dumps({"value": {"code": readfile(filename)}})
    res = requests.post("http://localhost:8080/init", data=body)
    return res.text

def run(body):
    #print("sending: %s" % body)
    return requests.post("http://localhost:8080/run", data='{"value": %s }' % body).text
 
usage = "usage: init <file>|run <json>"
if __name__ == '__main__':
    args = sys.argv[1:]
    if len(args) < 2: 
        print(usage)
    elif args[0].startswith("i"):
        print(init(args[1]))
    elif args[0].startswith("r"):
        print(run(args[1]))
    else: 
        print(usage)
