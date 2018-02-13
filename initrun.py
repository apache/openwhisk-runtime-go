import requests
import sys
import json
import six
import base64

def readfile(name, encoded=False):
    fp = open(name, "rb" if encoded else "r")
    data = fp.read()
    fp.close()
    if encoded:
        return base64.b64encode(data).decode('utf-8')
    return data

#filename="hello.py"
def code(filename):
    body = json.dumps({"value": {"code": readfile(filename)}})
    res = requests.post("http://localhost:8080/init", data=body)
    return res.text

def binary(filename):
    body = json.dumps({"value": {"binary": True, "code": readfile(filename, True)}})
    res = requests.post("http://localhost:8080/init", data=body)
    return res.text

def run(body):
    #print("sending: %s" % body)
    return requests.post("http://localhost:8080/run", data='{"value": %s }' % body).text
 
usage = "usage: code <file>|binary <file>|run <json>"
if __name__ == '__main__':
    args = sys.argv[1:]
    if len(args) < 2: 
        print(usage)
    elif args[0].startswith("c"):
        print(code(args[1]))
    elif args[0].startswith("b"):
        print(binary(args[1]))
    elif args[0].startswith("r"):
        print(run(args[1]))
    else: 
        print(usage)
