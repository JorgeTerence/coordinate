package main

import (
	"fmt"
	"net/http"
	"path"
	"os"
	qr "github.com/mdp/qrterminal/v3"
)

const (
	PORT         int32  = 8080
	INSTALL_PATH string = "/usr/share/coordinate"
)

var (
	pwd, host, addr = loadEnv()
	baseDir         = resolveBaseDir()
)

// TODO: config file for colors, filters, messages etc.
// TODO: Add favicon, logo and repo assets
func main() {
	http.HandleFunc("/", browse)

	programFiles := http.FileServer(http.Dir(path.Join(INSTALL_PATH, "web")))
	http.Handle("/static/", http.StripPrefix("/static/", programFiles))

	downloadFiles := http.FileServer(http.Dir(baseDir))
	http.Handle("/download/", http.StripPrefix("/download/", downloadFiles))

	http.HandleFunc("/tar/", downloadTar)

	url := fmt.Sprintf("http://%s:%d", addr, PORT)

	fmt.Printf("Serving from %s on %s\n", host, url)
	fmt.Printf("Base directory: %s\n\n", baseDir)
	qr.GenerateHalfBlock(url, qr.L, os.Stdout)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
