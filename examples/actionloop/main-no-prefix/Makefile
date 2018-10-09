OW_RUNTIME?=openwhisk/actionloop
OW_COMPILER?=openwhisk/actionloop-golang-v1.11
WSK?=wsk
MAIN=main
SRC=src/exec
NAME=golang-$(MAIN)-no-prefix
ZIP=$(MAIN).zip

deploy: test.done $(ZIP)
	$(WSK) action update test/$(NAME) $(ZIP) --main $(MAIN) --docker $(OW_RUNTIME)

$(ZIP): $(SRC)
	docker run -i $(OW_COMPILER) -compile $(MAIN) <$(SRC) >$(ZIP)

clean:
	wsk action delete test/$(NAME)
	rm $(MAIN) $(ZIP) test.done

test: deploy
	$(WSK) action invoke test/$(NAME) -r
	$(WSK) action invoke test/$(NAME) -P name.json -r

test.done:
	$(WSK) package update test
	touch test.done

.PHONY: deploy test clean
