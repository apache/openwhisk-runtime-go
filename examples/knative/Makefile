USR?=opennwhisk
BASE=actionloop-golang-v1.12
IMG=docker.io/$(USR)/$(BASE)
OUT=hello-go

.PHONY: push
push:
	docker tag $(BASE) $(IMG)
	docker login
	docker push $(IMG)

.PHONY: test-compile
test-compile:
	docker run -i $(IMG) \
	-compile main \
	 <../golang-main-single/main.go \
	 >exec.zip
	 unzip -l exec.zip

.PHONY: test-build
test-build:
	mkdir workspace 
	git clone https://gist.github.com/sciabarracom/a84ce9baa7d8eb7a769f80c61a618694 workspace/source
	docker run -i --entrypoint=/bin/bash \
	-v $(PWD)/workspace:/workspace \
	$(IMG) -c "/bin/proxy -compile main </workspace/source/src >/workspace/source/exec.zip"
	docker build -t $(OUT) workspace/source 

.PHONY: enter
enter:
	docker run \
	  -v $(PWD)/workspace:/workspace \
	  --entrypoint=/bin/bash \
	  -ti $(IMG)



