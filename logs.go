package main

import "log"

func throw(level Msg, msg interface{}) {
	switch level {
	case Get:
		log.Printf("\033[32mGET:\033[0m %s", msg)
	case Error:
		log.Printf("\033[31mERROR:\033[0m %s", msg)
	case Warn:
		log.Printf("\033[33mWARNING:\033[0m %s", msg)
	case Zip:
		log.Printf("\033[33mZIP:\033[0m %s", msg)
	}
}
