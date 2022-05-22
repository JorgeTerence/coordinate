package main

import (
	"archive/zip"
	"io"
	"os"
	"path"
)

func zipDir(source string, archive string) (string, error) {
	// Create archive file and file handler
	f, err := os.CreateTemp("", archive)
	if err != nil { return "", err }
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	// Recursively add all files and directories in the target directory
	return f.Name(), addDirToArchive(w, archive, source, "")
}

func addDirToArchive(w *zip.Writer, archive string, dirPath string, nestPath string) error {
	// Read list of files
	dir, err := os.ReadDir(dirPath)
	if err != nil { return err }

	for _, entry := range dir {
		// If the entry is a directory, call the cuntion recursively
		// `dirPath` is where the function tries to read files from
		// `nestPath` is the then directory structure inside the zip archive
		if entry.IsDir() {
			if err := addDirToArchive(w, archive, path.Join(dirPath, entry.Name()), path.Join(nestPath, entry.Name())); err != nil {
				return err
			}

			continue
		}

		// Else add regular file to archive
		file, err := os.Open(path.Join(dirPath, entry.Name()))
		if err != nil { return err }

		info, err := file.Stat()
		if err != nil { return err }

		header, err := zip.FileInfoHeader(info)
		if err != nil { return err }

		header.Method = zip.Deflate
		header.Name = path.Join(nestPath, header.Name)

		headerWriter, err := w.CreateHeader(header)
		if err != nil { return err }

		_, err = io.Copy(headerWriter, file)
		if err != nil { return err }
	}

	return nil
}
