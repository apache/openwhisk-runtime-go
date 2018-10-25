OW_USER?=openwhisk
OW_RUNTIME?=$(OW_USER)/actionloop
OW_COMPILER?=$(OW_USER)/actionloop-golang-v1.11
WSK?=wsk
MAIN=main
PACKAGE=test
SRCS=src/main/main.go
VENDORS=src/main/vendor
NAME=golang-$(MAIN)-vendor
BINZIP=$(MAIN)-bin.zip
SRCZIP=$(MAIN)-src.zip

deploy: package.done $(BINZIP)
	$(WSK) action update $(PACKAGE)/$(NAME) $(BINZIP) --main $(MAIN) --docker $(OW_RUNTIME)

devel: package.done $(SRCZIP)
	$(WSK) action update $(PACKAGE)/$(NAME) $(SRCZIP) --main $(MAIN) --docker $(OW_COMPILER)

src/%/vendor:
	cd $(@D) ; GOPATH=$$(realpath $$PWD/../..) dep ensure

$(BINZIP): $(SRCS) $(VENDORS)
	cd src ; zip - -r * | docker run -i $(OW_COMPILER) -compile $(MAIN) >../$(BINZIP)

$(SRCZIP): $(SRCS) $(VENDORS)
	cd src ; zip ../$(SRCZIP) -r *

clean:
	-$(WSK) action delete $(PACKAGE)/$(NAME)
	-rm $(BINZIP) $(SRCZIP) package.done test.json

clean_vendor:
	-rm -r $(VENDORS)

test: test.json
	$(WSK) action invoke test/$(NAME) -r
	$(WSK) action invoke test/$(NAME) -P test.json -r

test.json:
	echo '{ "name": "Mike" }' >test.json

package.done:
	$(WSK) package update $(PACKAGE)
	touch package.done

.PHONY: deploy devel test clean
