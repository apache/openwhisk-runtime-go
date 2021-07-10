/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package openwhisk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"
)

func startTestServer(compiler string) (*httptest.Server, string, *os.File) {
	// temporary workdir
	cur, _ := os.Getwd()
	dir, _ := ioutil.TempDir("", "action")
	file, _ := filepath.Abs("_test")
	os.Symlink(file, dir+"/_test")
	os.Chdir(dir)
	log.Printf(dir)
	// setup the server
	buf, _ := ioutil.TempFile("", "log")
	ap := NewActionProxy(dir, compiler, buf, buf)
	ts := httptest.NewServer(ap)
	log.Printf(ts.URL)
	doPost(ts.URL+"/init", `{value: {code: ""}}`)
	return ts, cur, buf
}

func stopTestServer(ts *httptest.Server, cur string, buf *os.File) {
	runtime.Gosched()
	// wait 2 seconds before declaring a test done
	time.Sleep(2 * time.Second)
	os.Chdir(cur)
	ts.Close()
	dump(buf)
}

func doPost(url string, message string) (string, int, error) {
	buf := bytes.NewBufferString(message)
	res, err := http.Post(url, "application/json", buf)
	if err != nil {
		return "", -1, err
	}
	defer res.Body.Close()
	resp, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", -1, err
	}
	return string(resp), res.StatusCode, nil
}

func doRun(ts *httptest.Server, message string) {
	if message == "" {
		message = `{"name":"Mike"}`
	}
	resp, status, err := doPost(ts.URL+"/run", `{ "value": `+message+`}`)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d %s", status, resp)
	}
	if !strings.HasSuffix(resp, "\n") {
		fmt.Println()
	}
}

func doInit(ts *httptest.Server, message string) {
	resp, status, err := doPost(ts.URL+"/init", message)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%d %s", status, resp)
	}
}

func initCode(file string, main string) string {
	dat, _ := ioutil.ReadFile(file)
	body := initBodyRequest{Code: string(dat)}
	if main != "" {
		body.Main = main
	}
	j, _ := json.Marshal(initRequest{Value: body})
	return string(j)
}

func initBytes(dat []byte, main string) string {
	enc := base64.StdEncoding.EncodeToString(dat)
	body := initBodyRequest{Binary: true, Code: enc}
	if main != "" {
		body.Main = main
	}
	j, _ := json.Marshal(initRequest{Value: body})
	return string(j)
}

func initBinary(file string, main string) string {
	dat, _ := ioutil.ReadFile(file)
	return initBytes(dat, main)
}

func abs(in string) string {
	out, _ := filepath.Abs(in)
	return out
}

func dump(file *os.File) {
	//file.Read()
	buf, _ := ioutil.ReadFile(file.Name())
	fmt.Print(string(buf))
	//fmt.Print(file.ReadAll())
	os.Remove(file.Name())
}

// printing output only if no errors
func sys(cli string, args ...string) {
	os.Chmod(cli, 0755)
	cmd := exec.Command(cli, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Print(err)
	} else {
		fmt.Print(string(out))
	}
}

// version printing output also when errors
func sys2(cli string, args ...string) {
	os.Chmod(cli, 0755)
	cmd := exec.Command(cli, args...)
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		log.Print(err)
	}
}

func exists(dir, filename string) error {
	path := fmt.Sprintf("%s/%d/%s", dir, highestDir(dir), filename)
	_, err := os.Stat(path)
	return err
}

func detectExecutable(dir, filename string) bool {
	path := fmt.Sprintf("%s/%d/%s", dir, highestDir(dir), filename)
	file, _ := ioutil.ReadFile(path)
	return IsExecutable(file, runtime.GOOS)
}

func waitabit() {
	time.Sleep(2000 * time.Millisecond)
}

func removeLineNr(out string) string {
	var re = regexp.MustCompile(`:\d+:\d+`)
	return re.ReplaceAllString(out, "::")
}
func TestMain(m *testing.M) {
	Debugging = false // enable debug of tests
	if !Debugging {
		// silence those annoying tests
		log.SetOutput(ioutil.Discard)
		// build support files
		sys("_test/build.sh")
		sys("_test/zips.sh")
	}

	// increase timeouts for init
	DefaultTimeoutStart = 1000 * time.Millisecond
	// build some test stuff
	// go ahead
	code := m.Run()
	os.Exit(code)
}
