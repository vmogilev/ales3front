package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func (app *AppConfig) GetRoutes() Routes {
	var routes Routes
	routes = Routes{
		Route{"001", "GET", "/", app.Index},
		Route{"002", "GET", "/s/term=", app.Search},
		Route{"003", "GET", "/s/term={term}", app.Search},
		Route{"003", "GET", "/" + app.DlPref + "/{s3path:.*}", app.Download},
	}
	return routes
}
