package main

import (
	"fmt"
	"net/http"
	//"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront/sign"
	"github.com/aws/aws-sdk-go/service/s3"
	//"github.com/nu7hatch/gouuid"
	"github.com/vmogilev/dlog"
)

type s3File struct {
	ContentLength int64
	ContentType   string
	LastModified  time.Time
}

// Page is used in the template
type Page struct {
	Title     string
	URLPath   string
	RawURL    string
	SignedURL string
	OK        bool
	Error     string
	Meta      *s3File
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
		dlog.Error.Printf("Failed to sign url, err: %s", err.Error())
		ok = false
		if c.debug {
			message = fmt.Sprintf("Failed to sign url, err: %s", err.Error())
		} else {
			message = "Failed to sign url, please notify site administrator!"
		}
	}
	return ok, message, signedURL
}

func (c *appContext) headS3File(key string) (bool, string, *s3File) {
	sess := session.New(&aws.Config{
		Region:      aws.String(c.region),
		Credentials: credentials.NewSharedCredentials("", c.cred),
	})
	svc := s3.New(sess)

	params := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	resp, err := svc.HeadObject(params)
	if err != nil {
		dlog.Error.Printf("couldn't get head of file: %s, %s", key, err)
		message := "File is not found!  Please check the URL"
		if c.debug {
			message = fmt.Sprintf("File is not found! Error: %s", err)
		}
		return false, message, &s3File{}
	}

	f := &s3File{
		ContentLength: *resp.ContentLength,
		ContentType:   *resp.ContentType,
		LastModified:  *resp.LastModified,
	}

	if c.debug {
		dlog.Trace.Println(resp)
		dlog.Trace.Println(f)
	}

	return true, "", f

}

func (c *appContext) dispatchHandler(w http.ResponseWriter, r *http.Request) {
	urlpath := r.URL.Path[0:]
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", c.bucket, urlpath)
}

func (c *appContext) cdnHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	t := r.FormValue("t")
	message := ""

	v := time.Now()
	ok := c.validateToken(t)
	respTime("validateToken", time.Since(v))

	if !ok {
		message = fmt.Sprintf("Download Token: %s is invalid", t)
	}

	//rawURL := path.Join(c.host, urlpath) // this strips :// to :/
	urlpath := r.URL.Path[len(c.cdn):]
	rawURL := c.host + urlpath
	signedURL := "#"
	meta := &s3File{}

	if ok {
		s := time.Now()
		ok, message, meta = c.headS3File(urlpath)
		respTime("headS3File", time.Since(s))
	}

	if ok {
		ok, message, signedURL = c.signURL(rawURL)
	}

	p := &Page{
		Title:     "Download",
		URLPath:   urlpath,
		RawURL:    rawURL,
		SignedURL: signedURL,
		OK:        ok,
		Error:     message,
		Meta:      meta,
	}
	c.renderTemplate(w, "download", p)
	respTime("cdnHandler", time.Since(start))
	//fmt.Fprintf(w, "<h1>%s</h1><div>urlpath: %s</div><div>rawURL: %s</div><div>signedURL: %s</div>", c.cdn, urlpath, rawURL, signedURL)

}
