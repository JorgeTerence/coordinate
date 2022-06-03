package main

import (
	"archive/tar"
	"os"
	"path"
)

func copyToTar(w *tar.Writer, source, nest string) error {
	info, err := os.Stat(source)
	if err != nil { return err }

	if info.IsDir() {
		dir, err := os.ReadDir(source)
		if err != nil { return err }

		for _, entry := range dir {
			if err := copyToTar(w, path.Join(source, entry.Name()), path.Join(nest, entry.Name())); err != nil {
				return err
			}
		}

		return nil
	}

	f, err := os.Open(source)
	if err != nil { return err }
	defer f.Close()

	content, err := os.ReadFile(source)
	if err != nil { return err }

	header := &tar.Header{
		Name: nest,
		Mode: int64(info.Mode()),
		Size: int64(len(string(content))),
	}

	if err := w.WriteHeader(header); err != nil {
		return err
	}

	_, err = w.Write(content)
	return err
}
