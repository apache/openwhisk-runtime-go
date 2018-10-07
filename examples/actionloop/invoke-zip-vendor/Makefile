MAIN=hello
OW_RUNTIME?=openwhisk/actionloop
OW_COMPILER?=openwhisk/actionloop-golang-v1.11
SRC=src/$(MAIN).go
NAME=golang-$(MAIN)-vendor
ZIP=$(MAIN).zip
WSK ?= wsk

deploy: test.done $(ZIP)
	$(WSK) action update test/$(NAME) $(ZIP) --main $(MAIN) --docker $(OW_RUNTIME) -p action test/golang-main-single

src/%/vendor:
	cd $(@D) ; dep ensure


$(ZIP): $(SRC) src/wskcli/vendor 
	cd src ; zip - -r * | docker run -i $(OW_COMPILER) -compile $(MAIN) >../$(ZIP)

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
