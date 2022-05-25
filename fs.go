package main

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
)

// FIXME: Final archive is broken - check if problem lies on back-end (here) or front-end (sending to browser)
func zipDir(source string) ([]byte, error) {
	// Create archive file and file handler
	f, err := os.CreateTemp("", fmt.Sprint(rand.Int()))
	if err != nil { return nil, err }
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	// Recursively add all files and directories in the target directory
	if err = copyToZip(w, source, ""); err != nil {
		return nil, err 
	}

	// Read archive's content and delete the file
	archive, err := os.ReadFile(f.Name())
	if err != nil { return nil, err }

	if err := os.Remove(f.Name()); err != nil {
		return nil, err
	}

	return archive, nil
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
	} 
	
	// Else (it's a file), add it to the archive
	file, err := os.Open(source)
	if err != nil { return err }
	defer file.Close()

	header, err := zip.FileInfoHeader(info)
	if err != nil { return err }

	header.Method = zip.Deflate
	header.Name = path.Join(nest, header.Name)

	headerWriter, err := w.CreateHeader(header)
	if err != nil { return err }

	_, err = io.Copy(headerWriter, file)
	return err
}
