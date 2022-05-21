package main

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

func zipDir(source string, target string, pageData BaseData) error {
	f, err := os.Create(target)
	if err != nil { return err }
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	dir, err := os.ReadDir(source)
	if err != nil { return err }

	for _, entry := range dir {
		if err := addFileToZip(w, target, path.Join(source, entry.Name())); err != nil {
			return err
		}
	}

	return nil
}

func addFileToZip(w *zip.Writer, archive string, target string) error {
	file, err := os.Open(target)
	if err != nil { return err }

	info, err := file.Stat()
	if err != nil { return err }

	header, err := zip.FileInfoHeader(info)
	if err != nil { return err }

	header.Method = zip.Deflate

	headerWriter, err := w.CreateHeader(header)
	if err != nil { return err }

	_, err = io.Copy(headerWriter, file)
	if err != nil { return err }

	return nil
}
