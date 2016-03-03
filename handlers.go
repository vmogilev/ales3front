package main

import (
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
	"github.com/vmogilev/dlog"
)

func (c *appContext) dispatchHandler(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path[0:]
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", c.bucket, urlpath)
}

func (c *appContext) cdnHandler(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path[len(c.cdn):]

	signer := sign.NewURLSigner(c.keyID, c.privKey)
	rawURL := path.Join(c.host, urlpath)

	signedURL, err := signer.Sign(rawURL, time.Now().Add(1*time.Hour))
	if err != nil {
		dlog.Error.Printf("Failed to sign url, err: %s\n", err.Error())
		return
	}
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div><div>%s</div>", c.cdn, urlpath, signedURL)

}
