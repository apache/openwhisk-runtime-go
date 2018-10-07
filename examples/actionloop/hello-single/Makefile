MAIN=hello
OW_RUNTIME?=openwhisk/actionloop
OW_COMPILER=openwhisk/actionloop-golang-v1.11
SRC=$(MAIN).go
NAME=golang-$(MAIN)-single
ZIP=$(MAIN).zip
WSK?=wsk

deploy: test.done $(ZIP)
	$(WSK) action update test/$(NAME) $(ZIP) --main $(MAIN) --docker $(OW_RUNTIME)

$(ZIP): $(SRC)
	docker run -i $(OW_COMPILER) -compile $(MAIN) <$(SRC) >$(ZIP)

clean:
	wsk action delete test/$(NAME)
	rm $(ZIP) test.done

test: deploy
	$(WSK) action invoke test/$(NAME) -r
	$(WSK) action invoke test/$(NAME) -P name.json -r

test.done:
	$(WSK) package update test
	touch test.done

.PHONY: deploy test
