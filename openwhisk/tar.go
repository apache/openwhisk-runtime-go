package openwhisk

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func openTar(src []byte) (*tar.Reader, error) {
	// Create a new bytes.Reader from the input byte slice
	reader := bytes.NewReader(src)

	// Create a new gzip.Reader from the bytes.Reader
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// Create a new tar.Reader from the gzip.Reader
	tarReader := tar.NewReader(gzipReader)

	return tarReader, nil
}

func UnTar(src []byte, dest string) error {
	r, err := openTar(src)
	if err != nil {
		return err
	}
	Debug("open Tar")
	os.MkdirAll(dest, 0755)
	for {
		header, err := r.Next()
		fmt.Println("Err", err)
		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dest, header.Name)
		// isLink := header.FileInfo().Mode()&os.ModeSymlink == os.ModeSymlink

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, r); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
