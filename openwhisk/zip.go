package openwhisk

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Unzip extracts file and directories in the given destination folder
func Unzip(src []byte, dest string) error {
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

// Zip creates a zip with one file named "exec"
func Zip(file string) ([]byte, error) {
	// read the file, zip it and write it to stdout
	buf := new(bytes.Buffer)
	zwr := zip.NewWriter(buf)
	zf, err := zwr.Create("exec")
	if err != nil {
		return nil, err
	}
	filedata, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	_, err = zf.Write(filedata)
	if err != nil {
		return nil, err
	}
	err = zwr.Flush()
	if err != nil {
		return nil, err
	}
	err = zwr.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
