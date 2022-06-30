package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"math"
	"net"
	"os"
	"path"
	"strings"
)

func getIPAddr() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}

		for _, addr := range addrs {
			ip := addr.String()

			isIPv6 := strings.Contains(ip, ":")
			isLocalhost := strings.Contains(ip, "127.0.0.1")

			if !isIPv6 && !isLocalhost {
				return strings.Split(ip, "/")[0], nil
			}
		}
	}

	return "", errors.New("FAILED TO FIND IP ADDRESS")
}

func loadEnv() (pwd string, host string, addr string) {
	var err error

	pwd, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	host, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	addr, err = getIPAddr()
	if err != nil {
		log.Fatal(err)
	}

	return
}

func loadTmpl(tmplName string) *template.Template {
	tmpl, err := template.ParseFS(assets, "web/base.html", fmt.Sprintf("web/%s.html", tmplName))
	if err != nil {
		log.Fatal(err)
	}

	return tmpl
}

func loadBaseData(url string) BaseData {
	return BaseData{
		Host:  host,
		Addr:  addr,
		Path:  url,
		Split: strings.Split(url, "/")[1:],
		Config: config,

		Join: path.Join,
		Size: fileSize,
	}
}

func resolveBaseDir() string {
	if len(os.Args) <= 1 || os.Args[1] == "." {
		return pwd
	}

	if path.IsAbs(os.Args[1]) {
		return path.Clean(os.Args[1])
	}

	return path.Join(pwd, os.Args[1])
}

func fileSize(n int64) string {
	if n == 0 {
		return "0B"
	}

	units := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	// Getting log base 1000 of n
	exp := math.Floor(math.Log(float64(n)) / math.Log(1000))

	return fmt.Sprintf("%.1f%s", float64(n)/math.Pow(1000, exp), units[int(exp)])
}

func record(level, msg string) {
	switch strings.ToUpper(level) {
	case "GET":
		log.Printf("\033[32mGET:\033[0m %s", msg)
	case "ERROR":
		log.Printf("\033[31mERROR:\033[0m %s", msg)
	case "WARN":
		log.Printf("\033[33mWARNING:\033[0m %s", msg)
	case "ZIP":
		log.Printf("\033[33mZIP:\033[0m %s", msg)
	}
}
