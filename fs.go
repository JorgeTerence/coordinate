package main

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"io"
	"os"
	"path"
	"strings"
)

// Creates temporary file for a .zip or .tar.gz archive of `source`
// Remember to remove file afterwards: `os.Remove(archivePath)`
func createTempArchive(source, format string) (string, error) {

	// Create archive file and file handler
	f, err := os.CreateTemp("", path.Base(source))
	if err != nil { return "", err }
	defer f.Close()
	
	switch strings.Trim(strings.ToLower(format), " ") {
	case "zip":
		w := zip.NewWriter(f)
		defer w.Close()

		// Recursively add all files and directories in the target directory
		if err = copyToZip(w, source, ""); err != nil {
			return "", err
		}
	
	case "tar":
		w := tar.NewWriter(f)
		defer w.Close()

		if err := copyToTar(w, source, ""); err != nil {
			return "", err
		}

	default: 
		return "", errors.New("`format` MUST BE EITHER \"zip\" or \"tar\"")
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

