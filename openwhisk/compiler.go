package openwhisk

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/h2non/filetype"
)

// check if the code is already compiled
func isCompiled(file string) bool {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return false
	}
	kind, err := filetype.Match(buf)
	if err != nil {
		log.Println(err)
		return false
	}
	if kind.Extension == "elf" {
		return true
	}
	return false
}

// CompileAction will compile an anction in source format invoking a compiler
func (ap *ActionProxy) CompileAction(file string, main string) error {
	if ap.compiler == "" {
		return fmt.Errorf("No compiler defined")
	}
	log.Printf("compiling: %s %s %s", ap.compiler, file, main)
	cmd := exec.Command(ap.compiler, file, main)
	out, err := cmd.CombinedOutput()
	log.Print(string(out))
	return err
}
