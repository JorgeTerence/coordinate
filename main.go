package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	qr "github.com/mdp/qrterminal/v3"
)

const PORT int32 = 8080

var (
	//go:embed web
	assets          embed.FS
	pwd, host, addr = loadEnv()
	config          = loadConfig()
	source          = resolveBaseDir()
)

// TODO: Support for audio, pdf and binaries
// TODO: Upload functionality
func main() {
	http.HandleFunc("/", browse)

	programFiles := http.FileServer(http.FS(assets))
	http.Handle("/static/", http.StripPrefix("/static/", programFiles))

	downloadFiles := http.FileServer(http.Dir(source))
	http.Handle("/download/", http.StripPrefix("/download/", downloadFiles))

	http.HandleFunc("/zip/", downloadZip)

	url := fmt.Sprintf("http://%s:%d", addr, PORT)

	fmt.Printf("Serving from %s on %s\n", host, url)
	fmt.Printf("Base directory: %s\n\n", source)

	if config.QrCode {
		qr.GenerateHalfBlock(url, qr.L, os.Stdout)
	}

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
