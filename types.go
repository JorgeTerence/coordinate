package main

import "io/fs"

type BaseData struct {
	Host      string
	Addr      string
	Path      string
	SplitPath []string
	Util      UtilFuncs
}

type UtilFuncs struct {
	PathJoin    func(...string) string
	ArrContains func([]string, string) bool
	Arr         func(...string) []string
}

type DirData struct {
	Base    BaseData
	Entries []fs.DirEntry
}

type FileData struct {
	Base    BaseData
	Content string
	Name    string
	Ext     string
}

type ErrorData struct {
	Base BaseData
	Err  error
}
