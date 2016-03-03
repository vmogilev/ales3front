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

func (c *appContext) dispatchHandler(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path[0:]
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", c.bucket, urlpath)
}

func (c *appContext) cdnHandler(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path[len(c.cdn):]

	signer := sign.NewURLSigner(c.keyID, c.privKey)
	//rawURL := path.Join(c.host, urlpath)
	rawURL := c.host + urlpath

	signedURL, err := signer.Sign(rawURL, time.Now().Add(time.Duration(c.expHours)*time.Hour))
	if err != nil {
		dlog.Error.Printf("Failed to sign url, err: %s\n", err.Error())
	}
	fmt.Fprintf(w, "<h1>%s</h1><div>urlpath: %s</div><div>rawURL: %s</div><div>signedURL: %s</div>", c.cdn, urlpath, rawURL, signedURL)

}
