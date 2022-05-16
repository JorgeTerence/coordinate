package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"github.com/pkg/browser"
)

const (
	PORT         int32  = 8080
	INSTALL_PATH string = "/home/jorge/Desktop/coordinate"
)

func main() {
	pwd, host, addr := getEnv()

	dirTmpl := loadTmpl("directory", INSTALL_PATH)
	fileTmpl := loadTmpl("file", INSTALL_PATH)
	errTmpl := loadTmpl("error", INSTALL_PATH)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	programFiles := http.FileServer(http.Dir(path.Join(INSTALL_PATH, "web")))
	http.Handle("/static/", http.StripPrefix("/static/", programFiles))

	downloadFiles := http.FileServer(http.Dir(pwd))
	http.Handle("/download/", http.StripPrefix("/download/", downloadFiles))

	fmt.Printf("Serving from %s on http://%s:%d\n", host, addr, PORT)
	browser.OpenURL(fmt.Sprintf("http://localhost:%d", PORT))
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
