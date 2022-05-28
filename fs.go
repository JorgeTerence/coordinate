package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path"
)

type Archive int

const (
	ZIP Archive = iota
	TAR
)

// Creates temporary file for a .zip or .tar.gz archive of `source`
// Remember to remove file afterwards: `os.Remove(archivePath)`
func createTempArchive(source string, format Archive) (string, error) {

	// Create archive file and file handler
	f, err := os.CreateTemp("", path.Base(source))
	if err != nil { return "", err }
	defer f.Close()
	
	switch format {
	case Archive(ZIP):
		w := zip.NewWriter(f)
		defer w.Close()

		// TODO: Implement both as a single function (pattern matching)
		// TODO: Better logic for this
		if err = copyToZip(w, source, path.Base(source)); err != nil {
			return "", err
		}
	
	case Archive(TAR):
		w := tar.NewWriter(f)
		defer w.Close()

		if err := copyToTar(w, source, path.Base(source)); err != nil {
			return "", err
		}
	}

	return f.Name(), err
}

func copyToZip(w *zip.Writer, source, nest string) error {
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
		Name: path.Join(nest, path.Base(f.Name())),
		Mode: int64(info.Mode()),
		Size: int64(len(string(content))),
	}

	if err := w.WriteHeader(header); err != nil {
		return err
	}

	_, err = w.Write(content)
	return err
}

func compress(source string) (string, error) {
	archive, err := os.CreateTemp("", path.Base(source))
	if err != nil { return "", err }
	defer archive.Close()

	w := gzip.NewWriter(archive)
	defer w.Close()

	f, err := os.Open(source)
	if err != nil { return "", err }

	_, err = io.Copy(w, f)
	return archive.Name(), err
}
