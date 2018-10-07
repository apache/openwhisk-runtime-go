MAIN=hello
RUSER ?= openwhisk
RUNTIME=$(RUSER)/actionloop
CUSER ?= openwhisk
COMPILER=$(CUSER)/actionloop-golang-v1.11
SRC=$(MAIN).go
NAME=golang-$(MAIN)-single
ZIP=$(MAIN).zip
WSK ?= wsk

deploy: test.done $(ZIP)
	$(WSK) action update test/$(NAME) $(ZIP) --main $(MAIN) --docker $(RUNTIME)

$(ZIP): $(SRC)
	docker run -i $(COMPILER) -compile $(MAIN) <$(SRC) >$(MAIN)
	zip $(ZIP) $(MAIN)

clean:
	wsk action delete test/$(NAME)
	rm $(MAIN) $(ZIP) test.done

test: deploy
	$(WSK) action invoke test/$(NAME) -r
	$(WSK) action invoke test/$(NAME) -P name.json -r

test.done:
	$(WSK) package update test
	touch test.done

.PHONY: deploy test
