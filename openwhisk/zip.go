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
	"log"
	"os"
	"path/filepath"
	"strings"
)

func openZip(src []byte) *zip.Reader {
	reader := bytes.NewReader(src)
	r, err := zip.NewReader(reader, int64(len(src)))
	if err != nil {
		return nil
	}
	return r
}

// UnzipOrSaveJar checks if is is a jar file looking if there is a META-INF folder in it
// if it is a jar file, save it as the file jarFile
// Otherwise unzip the files in the destination dir
func UnzipOrSaveJar(src []byte, dest string, jarFile string) error {
	r := openZip(src)
	if r == nil {
		return fmt.Errorf("not a zip file")
	}
	for _, f := range r.File {
		if f.Name == "META-INF/MANIFEST.MF" {
			ioutil.WriteFile(jarFile, src, 0644)
			return nil
		}
	}
	return Unzip(src, dest)
}

// Unzip extracts file and directories in the given destination folder
func Unzip(src []byte, dest string) error {
	r := openZip(src)
	os.MkdirAll(dest, 0755)
	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		defer rc.Close()
		if err != nil {
			return err
		}
		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			return os.MkdirAll(path, 0755)
		}
		err = os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		defer file.Close()
		if err != nil {
			return err
		}
		_, err = io.Copy(file, rc)
		return err
	}
	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

// Zip a directory
func Zip(dir string) ([]byte, error) {
	buf := new(bytes.Buffer)
	zwr := zip.NewWriter(buf)
	dir = filepath.Clean(dir)
	err := filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, dir)[1:]
		zipFile, err := zwr.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = zwr.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
