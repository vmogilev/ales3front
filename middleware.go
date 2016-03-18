package main

import (
	"net/http"
	"time"

	"github.com/vmogilev/dlog"
)

func logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		var client string
		client = r.Header.Get("X-Forwarded-For")
		if client == "" {
			client = r.RemoteAddr
		}

		dlog.Info.Printf(
			"%s\t%s\t%s\t%s\t%s",
			client,
			time.Since(start),
			r.Method,
			r.RequestURI,
			r.Referer(),
		)
	}

	return http.HandlerFunc(fn)
}

func recovery(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				c.simpleEvent("app", "panic")
				dlog.Error.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
