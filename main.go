package main

import (
	"fmt"
	"github.com/pkg/browser"
	"net/http"
	"path"
)

const (
	PORT         int32  = 8080
	INSTALL_PATH string = "/home/jorge/Desktop/coordinate"
)

func main() {
	pwd, host, addr := loadEnv()

	http.HandleFunc("/", browse)

	programFiles := http.FileServer(http.Dir(path.Join(INSTALL_PATH, "web")))
	http.Handle("/static/", http.StripPrefix("/static/", programFiles))

	downloadFiles := http.FileServer(http.Dir(pwd))
	http.Handle("/download/", http.StripPrefix("/download/", downloadFiles))

	fmt.Printf("Serving from %s on http://%s:%d\n", host, addr, PORT)
	browser.OpenURL(fmt.Sprintf("http://localhost:%d", PORT))
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
