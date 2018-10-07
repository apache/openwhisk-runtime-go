package wskcli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

func initEnv() {
	filename, _ := homedir.Expand("~/.wskprops")
	file, _ := os.Open(filename)
	defer file.Close()
	scan := bufio.NewScanner(file)
	for scan.Scan() {
		kv := strings.Split(scan.Text(), "=")
		switch kv[0] {
		case "APIHOST":
			os.Setenv("__OW_API_HOST", fmt.Sprintf("https://%s:443", kv[1]))
		case "AUTH":
			os.Setenv("__OW_API_KEY", kv[1])
		case "NAMESPACE":
			os.Setenv("__OW_NAMESPACE", kv[1])
		}
	}

}

func ExampleInvoke() {
	initEnv()
	payload := map[string]string{
		"hello": "world",
		"name":  "test",
	}
	res, err := Invoke("test/golang-main-single", payload)
	fmt.Printf("%s\n%v\n", res, err)
	// Output:
	// map[main-single:Hello, test!]
	// <nil>
}
