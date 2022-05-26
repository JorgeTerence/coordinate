package main

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

// Creates temporary file with a zip-archived copy of `source`
// Remember to remove file afterwards: `os.Remove(archivePath)`
func zipTmp(source string) (string, error) {

	// Create archive file and file handler
	f, err := os.CreateTemp("", path.Base(source))
	if err != nil { return "", err }
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	// Recursively add all files and directories in the target directory
	if err = copyToZip(w, source, ""); err != nil {
		return "", err
	}

	return f.Name(), err
}

func copyToZip(w *zip.Writer, source string, nest string) error {
	info, err := os.Stat(source)
	if err != nil { return err }

	// If the source is a directory, call the funtion recursively
	if info.IsDir() {
		dir, err := os.ReadDir(source)
		if err != nil { return err }

		for _, entry := range dir {
			if err := copyToZip(w, path.Join(source, entry.Name()), path.Join(nest, entry.Name())); err != nil {
				return err
			}
		}

		return nil
	}

	// Else (it's a file), add it to the archive
	file, err := os.Open(source)
	if err != nil { return err }
	defer file.Close()

	info, err = file.Stat()
	if err != nil { return err }

	header, err := zip.FileInfoHeader(info)
	if err != nil { return err }

	header.Method = zip.Deflate
	header.Name = path.Join(nest, header.Name)

	headerWriter, err := w.CreateHeader(header)
	if err != nil { return err }

	_, err = io.Copy(headerWriter, file)
	return err
}
