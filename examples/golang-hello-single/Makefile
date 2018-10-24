OW_USER?=openwhisk
OW_RUNTIME?=$(OW_USER)/actionloop
OW_COMPILER?=$(OW_USER)/actionloop-golang-v1.11
WSK?=wsk
MAIN=hello
PACKAGE=test
SRC=$(MAIN).go
NAME=golang-$(MAIN)-single
ZIP=$(MAIN).zip

deploy: package.done $(ZIP)
	$(WSK) action update $(PACKAGE)/$(NAME) $(ZIP) --main $(MAIN) --docker $(OW_RUNTIME)

devel: package.done $(SRC)
	$(WSK) action update $(PACKAGE)/$(NAME) $(SRC) --main $(MAIN) --docker $(OW_COMPILER)

$(ZIP): $(SRC)
	docker run -i $(OW_COMPILER) -compile $(MAIN) <$(SRC) >$(ZIP)

clean:
	-$(WSK) action delete $(PACKAGE)/$(NAME)
	-rm $(ZIP) package.done

test: test.json
	$(WSK) action invoke test/$(NAME) -r
	$(WSK) action invoke test/$(NAME) -P test.json -r

test.json:
	echo '{ "name": "Mike" }' >test.json

package.done:
	$(WSK) package update $(PACKAGE)
	touch package.done

.PHONY: deploy devel test clean
