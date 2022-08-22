package main

import (
	"html/template"
	"io"
)

func render(w io.Writer, t Msg, data interface{}) {
	f := []string{"directory.html", "file.html", "error.html"}[t]

	tmpl, err := template.ParseFS(assets, "web/base.html", "web/"+f)
	if err != nil {
		throw(Error, err.Error())
		return
	}

	tmpl.ExecuteTemplate(w, f, data)
}
