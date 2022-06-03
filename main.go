package main

import (
	"fmt"
	"net/http"
	"path"
)

const (
	PORT         int32  = 8080
	INSTALL_PATH string = "/usr/share/coordinate"
)

var (
	pwd, host, addr = loadEnv()
	baseDir         = resolveBaseDir()
)

func main() {
	http.HandleFunc("/", browse)

	programFiles := http.FileServer(http.Dir(path.Join(INSTALL_PATH, "web")))
	http.Handle("/static/", http.StripPrefix("/static/", programFiles))

	downloadFiles := http.FileServer(http.Dir(baseDir))
	http.Handle("/download/", http.StripPrefix("/download/", downloadFiles))

	http.HandleFunc("/tar/", downloadTar)

	fmt.Printf("Serving from %s on http://%s:%d\n", host, addr, PORT)
	fmt.Printf("Base directory: %s\n", baseDir)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
