package main

import (
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
		Util      UtilFuncs
	}

	UtilFuncs struct {
		PathJoin    func(...string) string
		ArrContains func([]string, string) bool
		Arr         func(...string) []string
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
	targetPath := path.Join(pwd, r.URL.Path)
	target, err := os.Stat(targetPath)
	pageData := BaseData{host, addr, r.URL.Path, strings.Split(r.URL.Path, "/")[1:], UtilFuncs{path.Join, contains, createArr}}

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

func upload(w http.ResponseWriter, r *http.Request) {
	
}
