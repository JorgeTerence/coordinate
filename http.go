package main

import (
	"archive/zip"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/samber/lo"
)

var (
	dirTmpl  = loadTmpl("directory")
	fileTmpl = loadTmpl("file")
	errTmpl  = loadTmpl("error")

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
		FileSize    func(int64) string
	}

	DirData struct {
		Base      BaseData
		Entries   []fs.DirEntry
		IsRoot    bool
		DirName string
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
	targetPath := path.Join(source, r.URL.Path)
	target, err := os.Stat(targetPath)
	pageData := loadBaseData(r.URL.Path)

	log.Printf("\033[32mGET:\033[0m %s", path.Clean(r.URL.Path))

	if err != nil {
		report(w, err)
		return
	}

	if target.IsDir() {
		dir, err := os.ReadDir(targetPath)

		if err != nil {
			report(w, err)
			return
		}

		filtered := lo.Filter(dir, func(f fs.DirEntry, _ int) bool { return !lo.Contains(ignored, f.Name()) })
		isAbs := targetPath == "/"
		dirName := path.Base(targetPath)
		
		if err := dirTmpl.ExecuteTemplate(w, "directory.html", DirData{pageData, filtered, isAbs, dirName}); err != nil {
			log.Fatal(err)
		}
	} else {
		file, err := os.ReadFile(targetPath)

		if err != nil {
			report(w, err)
			return
		}

		name := path.Base(targetPath)
		extention := path.Ext(targetPath)
		size := fileSize(target.Size())

		fileTmpl.ExecuteTemplate(w, "file.html", FileData{pageData, string(file), name, extention, size})
	}
}

func downloadZip(w http.ResponseWriter, r *http.Request) {
	target := strings.TrimPrefix(r.URL.Path, "/zip/")
	targetPath := path.Join(source, target)

	log.Printf("\033[33mZIP:\033[0m %s", path.Clean(target))

	w.Header().Set("Content-Type", "application/zip")

	zw := zip.NewWriter(w)
	defer zw.Close()

	if err := copyToZip(zw, targetPath, path.Base(targetPath)); err != nil {
		report(w, err)
		return
	}
}
