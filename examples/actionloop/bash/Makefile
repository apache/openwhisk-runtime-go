WSK ?= wsk
NAME=hello-bash
FILE=hello.sh

deploy: test.done
	$(WSK) action update test/$(NAME) $(FILE) --docker openwhisk/actionloop

test: deploy
	$(WSK) action invoke test/$(NAME) -r
	$(WSK) action invoke test/$(NAME) -P name.json -r

test.done:
	$(WSK) package update test
	touch test.done

clean:
	$(WSK) action delete test/$(NAME)