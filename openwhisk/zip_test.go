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
