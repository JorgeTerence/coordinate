package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"os"
)

func getLocalAddr(host string) (net.Addr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.New("FAILED TO FIND NETWORK INTERFACES")
	}

	addrs, err := ifaces[2].Addrs()
	if err != nil {
		return nil, errors.New("FAILED TO FIND NETWORK INTERFACES ADDRESSES")
	}

	return addrs[0], nil
}

func getEnv() (pwd string, host string, addr net.Addr) {
	var err error

	pwd, err = os.Getwd()
	if err != nil {
		log.Fatal("Failed to find working direcory")
	}

	host, err = os.Hostname()
	if err != nil {
		log.Fatal("Failed to read machine name")
	}

	addr, err = getLocalAddr(host)
	if err != nil {
		log.Fatal("Failed to find machine address in the local network")
	}

	return
}

func loadTmpl(tmplName string, path string) (tmpl *template.Template) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/web/%s.html", path, tmplName))
	if err != nil {
		log.Fatal(err)
	}

	return
}
