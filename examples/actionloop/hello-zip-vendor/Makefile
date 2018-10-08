OW_RUNTIME?=openwhisk/actionloop
OW_COMPILER?=openwhisk/actionloop-golang-v1.11
WSK?=wsk
MAIN=hello
SRC=src/$(MAIN).go
NAME=golang-$(MAIN)-zip-vendor
ZIP=$(MAIN).zip

deploy: test.done $(ZIP)
	$(WSK) action update test/$(NAME) $(ZIP) --main $(MAIN) --docker $(OW_RUNTIME)

$(ZIP): $(SRC) src/hello/*.go src/hello/vendor 
	cd src ; zip - -r * | docker run -i $(OW_COMPILER) -compile $(MAIN) >../$(ZIP)

test: deploy
	$(WSK) action invoke test/$(NAME) -r
	$(WSK) action invoke test/$(NAME) -P name.json -r

clean:
	wsk action delete test/$(NAME)
	rm $(MAIN) $(ZIP) test.done

src/%/vendor:
	cd $(@D) ; dep ensure

test.done:
	$(WSK) package update test
	touch test.done

.PHONY: deploy test clean
