package main

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/vmogilev/dlog"
)

func MountPoint(httpMount string) string {
	if httpMount == "/" {
		return ""
	} else {
		return httpMount
	}
}

func NewRouter(httpMount string, rootDir string, img string) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	mp := MountPoint(httpMount)
	routes := app.GetRoutes()

	dlog.Info.Println(routes)

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(mp + route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	router.PathPrefix(mp + "/img").Handler(http.StripPrefix(mp+"/img", Logger(http.FileServer(http.Dir(img)), "img")))
	router.PathPrefix(mp + "/").Handler(http.StripPrefix(mp+"/", Logger(http.FileServer(http.Dir(filepath.Join(rootDir, "static"))), "Static")))

	return router
}
