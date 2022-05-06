package main

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

const (
	PORT         int32  = 8080
	INSTALL_PATH string = "/home/jorge/Desktop/coordinate"
	// INSTALL_PATH string = "/usr/bin/coordinate"
)

type DirPageData struct {
	Files     []fs.DirEntry
	Host      string
	Addr      string
	SplitPath []string
	Path      string
	PathJoin  func(elem ...string) string
	Error     error
}

type FilePageData struct {
	Host      string
	Addr      string
	SplitPath []string
	FileData  string
}

func main() {
	pwd, host, address := getEnv()
	addr := address.String()[:len(address.String())-3]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// PROD: Load templates on the first executioon, not on every request
		dirTmpl := loadTmpl("directory", INSTALL_PATH)
		fileTmpl := loadTmpl("file", INSTALL_PATH)

		targetPath := path.Join(pwd, r.URL.Path)
		splitPath := strings.Split(r.URL.Path, "/")[1:]

		// Try to read a directory
		dir, err := os.ReadDir(targetPath)

		if err == nil {
			dirTmpl.Execute(w, DirPageData{dir, host, addr, splitPath, r.URL.Path, path.Join, nil})
		} else {
			// If that fails, try to read a file
			file, err := os.ReadFile(targetPath)

			if err == nil {
				fileTmpl.Execute(w, FilePageData{host, addr, splitPath, string(file)})
			} else {
				// If that fails, return an error message
				w.WriteHeader(404)
				dirTmpl.Execute(w, DirPageData{[]fs.DirEntry{}, host, addr, splitPath, "", path.Join, errors.New("DIRECTORY OR FILE NOT FOUND")})
			}
		}
	})

	// IDEA: Allow for custom templates and stylesheets inside '~/.config/coordinate'
	programFiles := http.FileServer(http.Dir(path.Join(INSTALL_PATH, "web")))
	http.Handle("/static/", http.StripPrefix("/static/", programFiles))

	fmt.Printf("Serving from %s on http://%s:%d\n", host, addr, PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
