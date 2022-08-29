package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	qr "github.com/mdp/qrterminal/v3"
)

type Msg = uint8

const (
	Directory Msg = iota
	File
	Error
	Warn
	Get
	Zip
)

const PORT int32 = 8080

var (
	//go:embed web
	assets          embed.FS
	pwd, host, addr = loadEnv()
	source          = resolveBase()

	urls = map[string]string{
		"static":   randUrl(1),
		"download": randUrl(2),
		"zip":      randUrl(3),
	}
)

func main() {
	programFiles := http.FileServer(http.FS(assets))
	http.Handle(urls["static"], http.StripPrefix(urls["static"], programFiles))

	downloadFiles := http.FileServer(http.Dir(source))
	http.Handle(urls["download"], http.StripPrefix(urls["download"], downloadFiles))

	http.HandleFunc(urls["zip"], downloadZip)
	http.HandleFunc("/", browse)

	url := fmt.Sprintf("http://%s:%d", addr, PORT)

	fmt.Printf("Serving from %s on %s\n", host, url)
	fmt.Printf("Base directory: %s\n\n", source)

	qr.GenerateHalfBlock(url, qr.L, os.Stdout)

	go http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	fmt.Scanln()
}
