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
	"strconv"
)

// higherDir will find the highest numeric name a sub directory has
// 0 if no numeric dir names found
func highestDir(dir string) int {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return 0
	}
	max := 0
	for _, file := range files {
		n, err := strconv.Atoi(file.Name())
		if err == nil {
			if n > max {
				max = n
			}
		}
	}
	return max
}

// ExtractAction accept a byte array and write it to a file
// it handles zip files extracting the content
// it stores in a new directory under ./action/XXX/suffix where x is incremented every time
// it returns the file if a file or the directory if it was a zip file
func (ap *ActionProxy) ExtractAction(buf *[]byte, suffix string) (string, error) {
	if buf == nil || len(*buf) == 0 {
		return "", fmt.Errorf("no file")
	}
	ap.currentDir++
	newDir := fmt.Sprintf("%s/%d/%s", ap.baseDir, ap.currentDir, suffix)
	os.MkdirAll(newDir, 0755)
	file := newDir + "/exec"
	if IsZip(*buf) {
		jar := os.Getenv("OW_SAVE_JAR")
		if jar != "" {
			jarFile := newDir + "/" + jar
			Debug("Extract Action, checking if it is a jar first")
			return jarFile, UnzipOrSaveJar(*buf, newDir, jarFile)
		}
		Debug("Extract Action, assuming a zip")
		return file, Unzip(*buf, newDir)

	} else if IsGz(*buf) {
		Debug("Extract Action, assuming a tar.gz")
		return file, UnTar(*buf, newDir)
	}
	return file, ioutil.WriteFile(file, *buf, 0755)
}
