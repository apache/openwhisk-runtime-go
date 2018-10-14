WSK?=wsk
OW_USER?=openwhisk
OW_RUNTIME?=$(OW_USER)/actionloop
NAME=bash-hello
PACKAGE=test
SRC=hello.sh

deploy: package.done
	$(WSK) action update $(PACKAGE)/$(NAME) $(SRC) --docker $(OW_RUNTIME)

test: 
	$(WSK) action invoke $(PACKAGE)/$(NAME) -r
	$(WSK) action invoke $(PACKAGE)/$(NAME) -P name.json -r

package.done:
	$(WSK) package update $(PACKAGE)
	touch package.done

clean:
	$(WSK) action delete $(PACKAGE)/$(NAME)

.PHONY: test deploy clean
