# prefix

USER=msciab     	# devel
#USER=openwhisk		# deploy

all: .PHONY
.PHONY: bash-hello golang-hello-single golang-main-single golang-main-standalone golang-main-package golang-hello-vendor 

bash-hello:
	OW_USER=$(USER) make deploy test -C bash-hello

golang-main-single:
	OW_USER=$(USER) make deploy test -C $@
	OW_USER=$(USER) make deploy test -C $@

golang-hello-single:
	OW_USER=$(USER) make deploy test -C $@
	OW_USER=$(USER) make deploy test -C $@

golang-main-standalone:
	OW_USER=$(USER) make deploy test -C $@
	OW_USER=$(USER) make deploy test -C $@

golang-main-package:
	OW_USER=$(USER) make deploy test -C $@
	OW_USER=$(USER) make deploy test -C $@

golang-hello-vendor:
	OW_USER=$(USER) make deploy test -C $@
	OW_USER=$(USER) make deploy test -C $@


