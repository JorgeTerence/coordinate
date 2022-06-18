package main

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

func copyToZip(w *zip.Writer, source, nest string) error {
	info, err := os.Stat(source)
	if err != nil { return err }

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

	file, err := os.Open(source)
	if err != nil { return err }
	defer file.Close()

	f, err := w.Create(nest)
	if err != nil { return err }

	_, err = io.Copy(f, file)
	return err
}
