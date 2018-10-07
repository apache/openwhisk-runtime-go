WSK?=wsk
OW_RUNTIME?=openwhisk/actionloop
NAME=hello-bash
FILE=hello.sh

deploy: test.done
	$(WSK) action update test/$(NAME) $(FILE) --docker $(OW_RUNTIME)

test: deploy
	$(WSK) action invoke test/$(NAME) -r
	$(WSK) action invoke test/$(NAME) -P name.json -r

test.done:
	$(WSK) package update test
	touch test.done

clean:
	$(WSK) action delete test/$(NAME)