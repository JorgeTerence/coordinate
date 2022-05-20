package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"os"
	"strings"
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

func loadTmpl(tmplName string, path string) (tmpl *template.Template) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/web/%s.html", path, tmplName))
	if err != nil { log.Fatal(err) }

	return
}

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value { return true }
	}

	return false
}

func createArr(args ...string) []string {
	return args
}
