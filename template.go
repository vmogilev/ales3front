package main

import (
	"html/template"
	//"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/vmogilev/dlog"
)

func (c *appContext) renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	defer respTime("renderTemplate")() // don't forget the extra parentheses

	t, err := template.ParseFiles(filepath.Join(c.htmlPath, tmpl) + ".html")
	if err != nil {
		dlog.Error.Printf("Failed to parse template, err: %s\n", err.Error())

		var message string
		if c.trace {
			message = err.Error()
		} else {
			message = "Failed to parse template: please notify server administrator"
		}

		http.Error(w, message, http.StatusInternalServerError)
		return
	}
	t.Execute(w, p)
}
