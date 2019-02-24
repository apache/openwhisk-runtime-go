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

import "fmt"

var zipFile = []byte{
	0x50, 0x4b, 0x03, 0x04, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x22, 0x60, 0x83, 0x4d, 0x00, 0x00,
}

var linuxFile = []byte{
	0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
}

var darwinFile = []byte{
	0xcf, 0xfa, 0xed, 0xfe, 0x07, 0x00, 0x00, 0x01, 0x03, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
}

var windowsFile = []byte{
	0x4d, 0x5a, 0x90, 0x00, 0x03, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00,
}

var shellFile = []byte("#!/bin/sh\necho hello\n")

func Example_filetype() {
	fmt.Printf("%t\n%t\n", IsElf(linuxFile), IsElf(zipFile))
	fmt.Printf("%t\n%t\n", IsMach64(darwinFile), IsMach64(zipFile))
	fmt.Printf("%t\n%t\n", IsExe(windowsFile), IsExe(zipFile))
	fmt.Printf("%t\n%t\n", IsZip(zipFile), IsExe(linuxFile))
	fmt.Printf("%t\n%t\n", IsExecutable(linuxFile, "linux"), IsExecutable(zipFile, "linux"))
	fmt.Printf("%t\n%t\n", IsExecutable(windowsFile, "windows"), IsExecutable(zipFile, "windows"))
	fmt.Printf("%t\n%t\n", IsExecutable(darwinFile, "darwin"), IsExecutable(zipFile, "darwin"))
	fmt.Printf("%t\n%t\n%t\n", IsExecutable(shellFile, "darwin"), IsExecutable(shellFile, "linux"), IsExecutable(shellFile, "windows"))
	// Output:
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// false
	// true
	// true
	// false

}
