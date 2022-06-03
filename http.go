package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"

	// "bytes"
	// "compress/gzip"
	// "io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	dirTmpl  = loadTmpl("directory", INSTALL_PATH)
	fileTmpl = loadTmpl("file", INSTALL_PATH)
	errTmpl  = loadTmpl("error", INSTALL_PATH)
)

type (
	BaseData struct {
		Host      string
		Addr      string
		Path      string
		SplitPath []string

		PathJoin    func(...string) string
		ArrContains func([]string, string) bool
		Arr         func(...string) []string
		Last        func([]string) string
	}

	DirData struct {
		Base    BaseData
		Entries []fs.DirEntry
	}

	FileData struct {
		Base    BaseData
		Content string
		Name    string
		Ext     string
	}

	ErrorData struct {
		Base BaseData
		Err  error
	}
)

func browse(w http.ResponseWriter, r *http.Request) {
	targetPath := path.Join(baseDir, r.URL.Path)
	target, err := os.Stat(targetPath)
	pageData := loadBaseData(r.URL.Path)

	log.Printf("\033[32mGET:\033[0m %s", path.Clean(r.URL.Path))

	if err != nil {
		errTmpl.Execute(w, ErrorData{pageData, err})
		log.Printf("\033[31mERROR:\033[0m %s", err)
		return
	}

	if target.IsDir() {
		dir, err := os.ReadDir(targetPath)
		dirTmpl.Execute(w, DirData{pageData, dir})

		if err != nil {
			log.Printf("\033[31mERROR:\033[0m %s", err)
			errTmpl.Execute(w, ErrorData{pageData, err})
		}
	} else {
		file, err := os.ReadFile(targetPath)
		fileTmpl.Execute(w, FileData{pageData, string(file), path.Base(targetPath), path.Ext(targetPath)})

		if err != nil {
			log.Printf("\033[31mERROR:\033[0m %s", err)
			errTmpl.Execute(w, ErrorData{pageData, err})
		}
	}
}

func downloadTar(w http.ResponseWriter, r *http.Request) {
	target := strings.TrimPrefix(r.URL.Path, "/tar/")
	pageData := loadBaseData(r.URL.Path)

	log.Printf("\033[33mTAR:\033[0m %s", path.Clean(target))

	w.Header().Set("Content-Type", "application/zip")

	var buf bytes.Buffer

	// Archive directory's contents to tarball
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	if err := copyToTar(tw, path.Join(baseDir, target), path.Base(target)); err != nil {
		log.Printf("\033[31mERROR:\033[0m %s", err)
		errTmpl.Execute(w, ErrorData{pageData, err})
		return
	}

	// Compress tarball using gzip
	// FIXME: `tar: Unexpected EOF in archive`
	zw := gzip.NewWriter(w)
	defer zw.Close()
	
	if _, err := zw.Write(buf.Bytes()); err != nil {
		log.Printf("\033[31mERROR:\033[0m %s", err)
		errTmpl.Execute(w, ErrorData{pageData, err})
		return
	}
}
