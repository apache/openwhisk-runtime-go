# Examples

This is a collection of examples. 
Each examples has a  Makefile with 4 targets:

- `make deploy` (or just make) deploys the action, precompiling it
- `make devel`  deploy the action in source format, leaving the compilation to the runtime
-  `make test` runs a simple test on the action; it should be already deployed
- `clean` remove intermediate files

THis is the list of the examples:

- [Simple Golang action](golang-main-single) main is `main.Main`
- [Simple Golang action](golang-hello-single) main is `main.Hello`
- [Golang action with a subpackage](golang-main-package) main is `main.Main` invoking a `hello.Hello`
- [Golang action with a subpackage and vendor folder](golang-hello-vendor) main is `main.Hello` invoking a `hello.Hello` using a dependency `github.com/sirupsen/logrus`
- [Standalone Golang Action](golang-main-standalone) main is `main.main`, implements the ActionLoop directly 
- [Simple Bash action](bash-hello) a simple bash script action implementing the ActionLoop directly
