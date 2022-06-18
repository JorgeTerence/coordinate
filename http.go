package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"github.com/samber/lo"
	"io"
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

	ignored = []string{"node_modules", "package-lock.json", "venv", "__pycache__"}
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
		Last        func([]string) (string, error)
		FileSize		func(int64) string
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
		Size    string
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

		if err != nil {
			log.Printf("\033[31mERROR:\033[0m %s", err)
			errTmpl.Execute(w, ErrorData{pageData, err})
			return
		}

		filtered := lo.Filter(dir, func(f fs.DirEntry, _ int) bool { return !lo.Contains(ignored, f.Name()) })

		dirTmpl.Execute(w, DirData{pageData, filtered})
	} else {
		file, err := os.ReadFile(targetPath)

		if err != nil {
			log.Printf("\033[31mERROR:\033[0m %s", err)
			errTmpl.Execute(w, ErrorData{pageData, err})
			return
		}

		name := path.Base(targetPath)
		extention := path.Ext(targetPath)
		size := fileSize(target.Size())

		fileTmpl.Execute(w, FileData{pageData, string(file), name, extention, size})
	}
}

func downloadTar(w http.ResponseWriter, r *http.Request) {
	target := strings.TrimPrefix(r.URL.Path, "/tar/")
	targetPath := path.Join(baseDir, target)
	pageData := loadBaseData(r.URL.Path)

	log.Printf("\033[33mTAR:\033[0m %s", path.Clean(target))

	w.Header().Set("Content-Type", "application/tar")

	var buf bytes.Buffer

	// Archive directory's contents to tarball
	tw := tar.NewWriter(&buf)
	defer tw.Close()

	if err := copyToTar(tw, targetPath, path.Base(targetPath)); err != nil {
		log.Printf("\033[31mERROR:\033[0m %s", err)
		errTmpl.Execute(w, ErrorData{pageData, err})
		return
	}

	// Compress tarball using gzip
	// FIXME: tar: Unexpected EOF in archive
	zw := gzip.NewWriter(w)
	defer zw.Close()

	if _, err := io.Copy(zw, &buf); err != nil {
		log.Printf("\033[31mERROR:\033[0m %s", err)
		errTmpl.Execute(w, ErrorData{pageData, err})
	}
}
