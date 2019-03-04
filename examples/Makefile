# prefix
IMAGES?=openwhisk     	

all: .PHONY

.PHONY: bash-hello golang-hello-single golang-main-single golang-main-standalone golang-main-package golang-main-vendor golang-hello-vendor 

bash-hello:
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-main-single:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-hello-single:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-main-standalone:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-main-package:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-main-vendor:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test

golang-hello-vendor:
	cd $@ && OW_USER=$(IMAGES) make clean devel test
	cd $@ && OW_USER=$(IMAGES) make clean deploy test


