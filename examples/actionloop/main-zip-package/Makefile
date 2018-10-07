RUNTIME?=openwhisk/actionloop
COMPILER?=openwhisk/actionloop-golang-v1.11
WSK?=wsk
MAIN=main
SRC=src/$(MAIN).go
NAME=golang-$(MAIN)-zip-package
ZIP=$(MAIN).zip

deploy: test.done $(ZIP)
	$(WSK) action update test/$(NAME) $(ZIP) --main $(MAIN) --docker $(RUNTIME)

$(ZIP): $(SRC)
	cd src ; zip - -r * | docker run -i $(COMPILER) -compile $(MAIN) >../$(MAIN)
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
