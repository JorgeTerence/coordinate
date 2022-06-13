package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
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

func loadBaseData(url string) BaseData {
	return BaseData {
		Host: host,
		Addr: addr,
		Path: url,
		SplitPath: strings.Split(url, "/")[1:],
		
		PathJoin: path.Join,
		ArrContains: contains,
		Arr: createArr,
		Last: lastOfArr,
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

func contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value { return true }
	}

	return false
}

func createArr(args ...string) []string {
	return args
}

func lastOfArr(arr []string) string {
	return arr[len(arr)-1]
}

func filter[K interface{}](arr []K, f func(K) bool) (res []K) {
	for _, v := range arr {
		if f(v) {
			res = append(res, v)
		}
	}
	return res
}
