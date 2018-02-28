package openwhisk

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/h2non/filetype"
	"github.com/stretchr/testify/assert"
)

func sys(cli string) {
	cmd := exec.Command(cli)
	out, err := cmd.CombinedOutput()
	if err == nil {
		fmt.Printf(">>>%s\n%s<<<\n", cli, string(out))
	} else {
		fmt.Println("KO")
		log.Print(err)
	}
}

func TestExtractActionTest_exec(t *testing.T) {
	//sys("pwd")
	// cleanup
	assert.Nil(t, os.RemoveAll("./action"))
	file, _ := ioutil.ReadFile("_test/exec")
	extractAction(&file)
	_, err := os.Stat("./action/exec")
	assert.Nil(t, err)
}

func detect(filename string) string {
	file, _ := ioutil.ReadFile(filename)
	kind, _ := filetype.Match(file)
	return kind.Extension
}
func TestExtractActionTest_exe(t *testing.T) {
	//sys("pwd")
	// cleanup
	assert.Nil(t, os.RemoveAll("./action"))
	// match  exe
	file, _ := ioutil.ReadFile("_test/exec")
	extractAction(&file)
	assert.Equal(t, detect("./action/exec"), "elf")
}

func TestExtractActionTest_zip(t *testing.T) {
	//sys("pwd")
	// cleanup
	assert.Nil(t, os.RemoveAll("./action"))
	// match  exe
	file, _ := ioutil.ReadFile("_test/exec.zip")
	extractAction(&file)
	assert.Equal(t, detect("./action/exec"), "elf")
	if _, err := os.Stat("./action/etc"); err != nil {
		t.Fail()
	}
	if _, err := os.Stat("./action/dir/etc"); err != nil {
		t.Fail()
	}
}

func TestHigherDir(t *testing.T) {
	assert.Equal(t, higherDir("./_test"), 0)
	assert.Equal(t, higherDir("./_test/first"), 3)
	assert.Equal(t, higherDir("./_test/second"), 17)
}
