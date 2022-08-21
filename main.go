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
	source          = resolveBase()
)

/*
TODO:
- Support for audio, pdf and binaries
- Use random string for URIs to avoid conflicts
*/
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

	qr.GenerateHalfBlock(url, qr.L, os.Stdout)
	
	go http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	fmt.Scan()
}
