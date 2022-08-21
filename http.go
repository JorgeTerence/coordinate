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
	dirTmpl  = loadTmpl("directory")
	fileTmpl = loadTmpl("file")
	errTmpl  = loadTmpl("error")

	contentTypes = map[string][]string{
		"image": {".png", ".jpg", ".jpeg", ".webp", ".svg", ".gif"},
		"video": {".mp4", ".mov", ".avi", ".wmf"},
		"audio": {".mp3", ".wav", ".ogg", ".aiff", ".opus"},
	}
)

type (
	BaseData struct {
		Host   string
		Addr   string
		Path   string
		Split  []string

		Join func(...string) string
		Size func(int64) string
	}

	DirData struct {
		Base    BaseData
		Entries []fs.DirEntry
		IsRoot  bool
		DirName string
	}

	FileData struct {
		Base    BaseData
		Content string
		Name    string
		Type    string
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

	record("GET", path.Clean(r.URL.Path))

	if err != nil {
		record("ERROR", err.Error())
		errTmpl.ExecuteTemplate(w, "error.html", ErrorData{pageData, err})
		return
	}

	if target.IsDir() {
		dir, err := os.ReadDir(targetPath)

		if err != nil {
			record("ERROR", err.Error())
			errTmpl.ExecuteTemplate(w, "error.html", ErrorData{pageData, err})
			return
		}
		
		isRoot := targetPath == source
		dirName := path.Base(targetPath)

		dirTmpl.ExecuteTemplate(w, "directory.html", DirData{pageData, dir, isRoot, dirName})
	} else {
		file, err := os.ReadFile(targetPath)

		if err != nil {
			record("ERROR", err.Error())
			errTmpl.ExecuteTemplate(w, "error.html", ErrorData{pageData, err})
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

		fileTmpl.ExecuteTemplate(w, "file.html", FileData{pageData, string(file), name, fileType, size})
	}
}

func downloadZip(w http.ResponseWriter, r *http.Request) {
	target := strings.TrimPrefix(r.URL.Path, "/zip/")
	targetPath := path.Join(source, target)

	record("ZIP", path.Clean(r.URL.Path))

	w.Header().Set("Content-Type", "application/zip")

	zw := zip.NewWriter(w)
	defer zw.Close()

	if err := copyToZip(zw, targetPath, path.Base(targetPath)); err != nil {
		record("ERROR", err.Error())
		errTmpl.ExecuteTemplate(w, "error.html", ErrorData{loadBaseData(""), err})
	}
}
