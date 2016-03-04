package main

import (
	"fmt"
	"net/http"
	//"path"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
	//"github.com/nu7hatch/gouuid"
	"github.com/vmogilev/dlog"
)

type Page struct {
	Title     string
	URLPath   string
	RawURL    string
	SignedURL string
	OK        bool
	Error     string
}

func (c *appContext) validateToken(t string) bool {
	if t == "" {
		return false
	}

	if t == c.root {
		return true
	}

	return false
}

func (c *appContext) signURL(rawURL string) (bool, string, string) {
	ok := true
	message := ""

	signer := sign.NewURLSigner(c.keyID, c.privKey)
	signedURL, err := signer.Sign(rawURL, time.Now().Add(time.Duration(c.expHours)*time.Hour))
	if err != nil {
		dlog.Error.Printf("Failed to sign url, err: %s\n", err.Error())
		ok = false
		if c.debug {
			message = fmt.Sprintf("Failed to sign url, err: %s\n", err.Error())
		} else {
			message = "Failed to sign url, please notify site administrator!"
		}
	}
	return ok, message, signedURL
}

func (c *appContext) dispatchHandler(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path[0:]
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", c.bucket, urlpath)
}

func (c *appContext) cdnHandler(w http.ResponseWriter, r *http.Request) {
	t := r.FormValue("t")
	message := ""
	ok := c.validateToken(t)
	if !ok {
		message = fmt.Sprintf("Download Token: %s is invalid", t)
	}

	//rawURL := path.Join(c.host, urlpath) // this strips :// to :/
	urlpath := r.URL.Path[len(c.cdn):]
	rawURL := c.host + urlpath

	signedURL := "#"
	if ok {
		ok, message, signedURL = c.signURL(rawURL)
	}

	p := &Page{
		Title:     c.cdn,
		URLPath:   urlpath,
		RawURL:    rawURL,
		SignedURL: signedURL,
		OK:        ok,
		Error:     message,
	}
	c.renderTemplate(w, "download", p)
	//fmt.Fprintf(w, "<h1>%s</h1><div>urlpath: %s</div><div>rawURL: %s</div><div>signedURL: %s</div>", c.cdn, urlpath, rawURL, signedURL)

}
