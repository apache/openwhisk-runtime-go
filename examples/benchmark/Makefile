IMG?=whisk/actionloop-golang-v1.11
IMG2?=whisk/actionloop

all: golang bash

test.lua:
	echo 'wrk.method = "POST"'>test.lua
	echo "wrk.body = '{\"value\":{\"name\":\"Mike\"}}'">>test.lua
	echo 'wrk.headers["Content-Type"] = "application/json"'>>test.lua

golang: test.lua
	docker run -d --name under-test --rm -p 8080:8080 $(IMG)
	invoke init main.go
	wrk -t1 -c1 -stest.lua http://localhost:8080/run
	docker kill under-test

bash: test.lua
	docker run -d --name under-test --rm -p 8080:8080 $(IMG2)
	invoke init main.sh
	wrk -t1 -c1 -stest.lua http://localhost:8080/run
	docker kill under-test

.PHONY: all golang bash
