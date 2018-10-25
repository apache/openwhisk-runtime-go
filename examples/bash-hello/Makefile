WSK?=wsk
OW_USER?=openwhisk
OW_RUNTIME?=$(OW_USER)/actionloop
NAME=bash-hello
PACKAGE=test
SRC=hello.sh

deploy: package.done
	$(WSK) action update $(PACKAGE)/$(NAME) $(SRC) --docker $(OW_RUNTIME)

devel: deploy

test: test.json
	$(WSK) action invoke $(PACKAGE)/$(NAME) -r
	$(WSK) action invoke $(PACKAGE)/$(NAME) -P test.json -r

clean:
	-$(WSK) action delete $(PACKAGE)/$(NAME)
	-rm package.done test.json

package.done:
	$(WSK) package update $(PACKAGE)
	touch package.done

test.json:
	echo '{ "name": "Mike" }' >test.json


.PHONY: test devel deploy clean
