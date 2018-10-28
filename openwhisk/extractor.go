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
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/h2non/filetype.v1"
)

func unzip(src []byte, dest string) error {
	reader := bytes.NewReader(src)
	r, err := zip.NewReader(reader, int64(len(src)))
	if err != nil {
		return err
	}

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			rc.Close()
		}()

		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				f.Close()
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}
	return nil
}

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
	kind, err := filetype.Match(*buf)
	if err != nil {
		return "", err
	}
	file := newDir + "/exec"
	if kind.Extension == "zip" {
		Debug("Extract Action, assuming a zip")
		return file, unzip(*buf, newDir)
	}
	return file, ioutil.WriteFile(file, *buf, 0755)
}
