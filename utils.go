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

	"github.com/samber/lo"
)

func getIPAddr() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil { log.Fatal(err) }

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
	if err != nil { log.Fatal(err) }

	host, err = os.Hostname()
	if err != nil { log.Fatal(err) }

	addr, err = getIPAddr()
	if err != nil { log.Fatal(err) }

	return
}

func loadTmpl(tmplName string, source string) *template.Template {
	tmpl, err := template.ParseFiles(path.Join(source, "web", tmplName+".html"))
	if err != nil { log.Fatal(err) }
	return tmpl
}

func loadBaseData(url string) BaseData {
	return BaseData{
		Host:      host,
		Addr:      addr,
		Path:      url,
		SplitPath: strings.Split(url, "/")[1:],

		PathJoin:    path.Join,
		ArrContains: lo.Contains[string],
		Arr:         arr[string],
		Last:        lo.Last[string],
		FileSize: 	 fileSize,
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

func arr[T any](args ...T) []T {
	return args
}

func fileSize(n int64) string {
	if n == 0 {
		return "0B"
	}

	units := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	// Getting log base 1000 of n
	exp := math.Floor(math.Log(float64(n)) / math.Log(1000))

	return fmt.Sprintf("%.1f%s", float64(n) / math.Pow(1000, exp), units[int(exp)])
}
