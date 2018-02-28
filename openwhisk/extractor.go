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
	"strconv"

	"github.com/h2non/filetype"
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
func higherDir(dir string) int {
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

var currentDir = higherDir("./action")

// extractAction accept a byte array write it to a file
func extractAction(buf *[]byte) error {
	if buf == nil || len(*buf) == 0 {
		return fmt.Errorf("no file")
	}
	currentDir++
	newDir := fmt.Sprintf("./action/%d", currentDir)
	os.MkdirAll(newDir, 0755)
	kind, err := filetype.Match(*buf)
	if err != nil {
		return err
	}
	log.Println(kind)
	if kind.Extension == "elf" {
		log.Println("Extract Action, assuming a binary")
		return ioutil.WriteFile(newDir+"/exec", *buf, 0755)
	}
	if kind.Extension == "zip" {
		log.Println("Extract Action, assuming a zip")
		return unzip(*buf, newDir)
	}
	return fmt.Errorf("unknown filetype %s", kind)
}
