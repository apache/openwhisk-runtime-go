OW_USER?=openwhisk
OW_RUNTIME?=$(OW_USER)/actionloop
OW_COMPILER?=$(OW_USER)/actionloop-golang-v1.11
WSK?=wsk
MAIN=hello
PACKAGE=test
SRCS=src/$(MAIN).go src/hello/hello.go
VENDORS=src/hello/vendor
NAME=golang-$(MAIN)-vendor
BINZIP=$(MAIN)-bin.zip
SRCZIP=$(MAIN)-src.zip

deploy: package.done $(BINZIP)
	$(WSK) action update $(PACKAGE)/$(NAME) $(BINZIP) --main $(MAIN) --docker $(OW_RUNTIME)

devel: package.done $(SRCZIP)
	$(WSK) action update $(PACKAGE)/$(NAME) $(SRCZIP) --main $(MAIN) --docker $(OW_RUNTIME)

src/%/vendor:
	cd $(@D) ; GOPATH=$$PWD/../.. dep ensure

$(BINZIP): $(SRCS) $(VENDORS)
	cd src ; zip - -r * | docker run -i $(OW_COMPILER) -compile $(MAIN) >../$(BINZIP)

$(SRCZIP): $(SRCS) $(VENDORS)
	cd src ; zip ../$(SRCZIP) -r *

clean:
	wsk action delete $(PACKAGE)/$(NAME)
	rm  $(BINZIP) $(SRCZIP) package.done

test: test.json
	$(WSK) action invoke test/$(NAME) -r
	$(WSK) action invoke test/$(NAME) -P test.json -r

test.json:
	echo '{ "name": "Mike" }' >test.json

package.done:
	$(WSK) package update $(PACKAGE)
	touch package.done

.PHONY: deploy devel test clean
