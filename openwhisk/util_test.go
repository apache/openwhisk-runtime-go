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
	"path/filepath"
	"runtime"
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
	ap := NewActionProxy(dir, compiler, buf)
	ts := httptest.NewServer(ap)
	log.Printf(ts.URL)
	doPost(ts.URL+"/init", `{value: {code: ""}}`)
	return ts, cur, buf
}

func stopTestServer(ts *httptest.Server, cur string, buf *os.File) {
	runtime.Gosched()
	time.Sleep(1 * time.Second)
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

func initBinary(file string, main string) string {
	dat, _ := ioutil.ReadFile(file)
	enc := base64.StdEncoding.EncodeToString(dat)
	body := initBodyRequest{Binary: true, Code: enc}
	if main != "" {
		body.Main = main
	}
	j, _ := json.Marshal(initRequest{Value: body})
	return string(j)
}

func abs(in string) string {
	out, _ := filepath.Abs(in)
	return out
}

func Example_json_init() {
	fmt.Println(initCode("", ""))
	fmt.Println(initCode("_test/etc", ""))
	fmt.Println(initCode("_test/etc", "world"))
	fmt.Println(initBinary("_test/etc", ""))
	fmt.Println(initBinary("_test/etc", "hello"))
	// Output:
	// {"value":{}}
	// {"value":{"code":"1\n"}}
	// {"value":{"code":"1\n","main":"world"}}
	// {"value":{"code":"MQo=","binary":true}}
	// {"value":{"code":"MQo=","binary":true,"main":"hello"}}
}

func dump(file *os.File) {
	//file.Read()
	buf, _ := ioutil.ReadFile(file.Name())
	fmt.Print(string(buf))
	//fmt.Print(file.ReadAll())
	os.Remove(file.Name())
}
