def main(args):
    name = "world"
    if "name" in args:
        name = args["name"]
    return {"python": "Hello, %s" % name }