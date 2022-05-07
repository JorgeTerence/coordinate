package main

import (
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

type PageData struct {
	Host      string
	Addr      string
	Path      string
	SplitPath []string
	PathJoin  func(elem ...string) string
}
type DirPageData struct {
	Base    PageData
	Entries []fs.DirEntry
}
type FilePageData struct {
	Base    PageData
	Content string
	Name    string
}
type ErrorPageData struct {
	Base PageData
	Err  error
}

func main() {
	pwd, host, address := getEnv()
	addr := address.String()[:len(address.String())-3]

	dirTmpl := loadTmpl("directory", INSTALL_PATH)
	fileTmpl := loadTmpl("file", INSTALL_PATH)
	errTmpl := loadTmpl("error", INSTALL_PATH)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		targetPath := path.Join(pwd, r.URL.Path)
		target, err := os.Stat(targetPath)
		pageData := PageData{host, addr, r.URL.Path, strings.Split(r.URL.Path, "/")[1:], path.Join}

		if err != nil {
			errTmpl.Execute(w, ErrorPageData{pageData, err})
			return
		}

		if target.IsDir() {
			dir, err := os.ReadDir(targetPath)
			dirTmpl.Execute(w, DirPageData{pageData, dir})

			if err != nil {
				errTmpl.Execute(w, ErrorPageData{pageData, err})
			}
		} else {
			file, err := os.ReadFile(targetPath)
			fileTmpl.Execute(w, FilePageData{pageData, string(file), path.Base(targetPath)})
			
			if err != nil {
				errTmpl.Execute(w, ErrorPageData{pageData, err})
			}
		}
	})

	// IDEA: Allow for custom templates and stylesheets inside '~/.config/coordinate'
	programFiles := http.FileServer(http.Dir(path.Join(INSTALL_PATH, "web")))
	http.Handle("/static/", http.StripPrefix("/static/", programFiles))

	downloadFiles := http.FileServer(http.Dir(pwd))
	http.Handle("/download/", http.StripPrefix("/download/", downloadFiles))

	fmt.Printf("Serving from %s on http://%s:%d\n", host, addr, PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
