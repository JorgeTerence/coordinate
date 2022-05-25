package main

import (
	"fmt"
	"io/fs"
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

	if err != nil {
		errTmpl.Execute(w, ErrorData{pageData, err})
		return
	}

	if target.IsDir() {
		dir, err := os.ReadDir(targetPath)
		dirTmpl.Execute(w, DirData{pageData, dir})

		if err != nil {
			errTmpl.Execute(w, ErrorData{pageData, err})
		}
	} else {
		file, err := os.ReadFile(targetPath)
		fileTmpl.Execute(w, FileData{pageData, string(file), path.Base(targetPath), path.Ext(targetPath)})

		if err != nil {
			errTmpl.Execute(w, ErrorData{pageData, err})
		}
	}
}

// TODO: Better separation of concerns
func downloadZip(w http.ResponseWriter, r *http.Request) {
	dirPath := path.Join(baseDir, strings.TrimPrefix(r.URL.Path, "/zip/"))
	archiveName := fmt.Sprintf("%s.zip", path.Base(dirPath))

	pageData := loadBaseData(r.URL.Path)

	archive, err := zipDir(dirPath, archiveName)

	if err != nil {
		errTmpl.Execute(w, ErrorData{pageData, err})
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Write(archive)
}

func downloadTar(w http.ResponseWriter, r *http.Request) {

}
