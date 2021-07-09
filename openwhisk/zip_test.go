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
	"fmt"
	"io/ioutil"
	"os"
)

func Example_zip() {
	os.RemoveAll("./action/unzip")
	os.Mkdir("./action/unzip", 0755)
	buf, err := Zip("_test/pysample")
	fmt.Println(err)
	err = UnzipOrSaveJar(buf, "./action/unzip", "./action/unzip/exec.jar")
	sys("_test/find.sh", "./action/unzip")
	fmt.Println(err)
	// Output:
	// <nil>
	// ./action/unzip
	// ./action/unzip/exec
	// ./action/unzip/lib
	// ./action/unzip/lib/action
	// ./action/unzip/lib/action/__init__.py
	// ./action/unzip/lib/action/main.py
	// ./action/unzip/lib/exec.py
	// <nil>
}
func Example_jar() {
	os.RemoveAll("./action/unzip")
	os.Mkdir("./action/unzip", 0755)
	buf, err := Zip("_test/jar")
	fmt.Println(err)
	err = UnzipOrSaveJar(buf, "./action/unzip", "./action/unzip/exec.jar")
	sys("_test/find.sh", "./action/unzip")
	fmt.Println(err)
	// Output:
	// <nil>
	// ./action/unzip
	// ./action/unzip/exec.jar
	// <nil>
}

func Example_venv() {
	os.RemoveAll("./action/unzip")
	os.Mkdir("./action/unzip", 0755)
	buf, err := Zip("_test/venv")
	fmt.Println(1, err)
	err = ioutil.WriteFile("/tmp/appo.zip", buf, 0644)
	fmt.Println(2, err)
	err = UnzipOrSaveJar(buf, "./action/unzip", "./action/unzip/exec.jar")
	sys("bash", "-c", "cd action/unzip/bin && find . -type l -name python && rm ./python")
	sys2("bash", "-c", "diff -qr _test/venv action/unzip 2>/dev/null")
	fmt.Println(3, err)
	// Output:
	// 1 <nil>
	// 2 <nil>
	// ./python
	// Only in _test/venv/bin: python
	// 3 <nil>

}
