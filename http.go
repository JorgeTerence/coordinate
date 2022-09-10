package main

import (
	"archive/zip"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/samber/lo"
)

var (
	contentTypes = map[string][]string{
		"image": {".png", ".jpg", ".jpeg", ".webp", ".svg", ".gif"},
		"video": {".mp4", ".mov", ".avi", ".wmf"},
		"audio": {".mp3", ".wav", ".ogg", ".aiff", ".opus"},
	}
)

type (
	BaseData struct {
		Host  string
		Addr  string
		Path  string
		Split []string
		Urls  map[string]string
		Join  func(...string) string
		Size  func(int64) string
	}

	DirData struct {
		Base    BaseData
		Entries []fs.DirEntry
		IsRoot  bool
		DirName string
	}

	FileData struct {
		Base       BaseData
		Content    string
		Name       string
		Type       string
		Size       string
		Executable bool
	}

	ErrorData struct {
		Base  BaseData
		Error error
	}
)

func browse(w http.ResponseWriter, r *http.Request) {
	target, err := os.Stat(path.Join(source, r.URL.Path))

	throw(Get, path.Clean(r.URL.Path))

	if err != nil {
		throw(Error, err)
		render(w, Error, ErrorData{loadBaseData(""), err})
		return
	}

	if target.IsDir() {
		dirView(w, r)
	} else {
		fileView(w, r)
	}
}

func fileView(w http.ResponseWriter, r *http.Request) {
	targetPath := path.Join(source, r.URL.Path)
	target, err := os.Stat(targetPath)
	pageData := loadBaseData(r.URL.Path)

	if err != nil {
		throw(Error, err)
		render(w, Error, ErrorData{pageData, err})
		return
	}

	file, err := os.ReadFile(targetPath)

	if err != nil {
		throw(Error, err)
		render(w, Error, ErrorData{pageData, err})
		return
	}

	name := path.Base(targetPath)
	size := fileSize(target.Size())

	var fileType string

	for contentType, ext := range contentTypes {
		if lo.Contains(ext, path.Ext(targetPath)) {
			fileType = contentType
			break
		}
	}

	render(w, File, FileData{pageData, string(file), name, fileType, size, target.Mode().Perm()&0111 != 0})
}

func dirView(w http.ResponseWriter, r *http.Request) {
	targetPath := path.Join(source, r.URL.Path)
	pageData := loadBaseData(r.URL.Path)
	dir, err := os.ReadDir(targetPath)

	if err != nil {
		throw(Error, err)
		render(w, Error, ErrorData{pageData, err})
		return
	}

	isRoot := targetPath == source
	dirName := path.Base(targetPath)

	render(w, Directory, DirData{pageData, dir, isRoot, dirName})
}

func downloadZip(w http.ResponseWriter, r *http.Request) {
	target := strings.TrimPrefix(r.URL.Path, urls["zip"])
	targetPath := path.Join(source, target)

	throw(Zip, path.Clean(r.URL.Path))

	w.Header().Set("Content-Type", "application/zip")

	zw := zip.NewWriter(w)
	defer zw.Close()

	if err := copyToZip(zw, targetPath, path.Base(targetPath)); err != nil {
		throw(Error, err)
		render(w, Error, ErrorData{loadBaseData(""), err})
	}
}
